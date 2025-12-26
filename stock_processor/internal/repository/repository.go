package repository

import (
	"context"
	"fmt"
	"stock-processor/config"
	"stock-processor/internal/model"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type repository struct {
 rdb *redis.Client
}

type Repository interface{
	 SetCache(ctx context.Context, stock model.StockPrice) error
   AddToHistory(ctx context.Context, stock model.StockPrice) error
	 GetCache(ctx context.Context,symbol string)(float64, error)
	GetAll(ctx context.Context) (map[string]float64,error) 
	 GetHistory(ctx context.Context,symbol string,limit int) ([]float64, error)
}

func NewPriceRepository(rdb *redis.Client) Repository {

	return  &repository{
		rdb: rdb,
	}
}

func (r *repository) SetCache(ctx context.Context,stock  model.StockPrice)error {

 return  r.rdb.Set(ctx,stock.Symbol,stock.Price,0).Err()
}

func (r *repository) AddToHistory(ctx context.Context,stock  model.StockPrice)error  {
	key:=stock.Symbol+":history"
	pipe:=r.rdb.Pipeline()
	pipe.ZAdd(ctx,key,redis.Z{
		Score: float64(stock.Timestamp.Unix()),
		Member: stock.Price,
	})
	 pipe.ZRemRangeByRank(ctx, key, 0, -101)
	  _, err := pipe.Exec(ctx)
    return err
}

func (r *repository) GetCache(ctx context.Context,symbol string)(float64, error)  {

return	r.rdb.Get(ctx,symbol).Float64()
	
}


func (r *repository) GetAll(ctx context.Context) (map[string]float64,error)  {
	prices:=make(map[string]float64)
	pipe:=r.rdb.Pipeline()
	cmds:=make(map[string]*redis.StringCmd)
	for _, symbol := range config.Symbols {
		cmds[symbol]=pipe.Get(ctx,symbol)
	}
	_,err:=pipe.Exec(ctx)
	if err!=nil {
		return nil,err
	}
	for symbol, cmd := range cmds {
	if price,err:=	cmd.Float64();err==nil {
		prices[symbol]=price
	}
		
	}
	return prices,nil
}

func (r *repository) GetHistory(ctx context.Context,symbol string,limit int) ([]float64, error)  {
	 key := symbol + ":history"
values, err :=	r.rdb.ZRevRange(ctx,key,0,int64(limit-1)).Result()
	  if err != nil {
        return nil, err
  }
	prices:=make([]float64,len(values))
	for i, val := range values {
	price, parseErr := strconv.ParseFloat(val, 64)
    if parseErr != nil {
        return nil, fmt.Errorf("failed to parse price at index %d: %w", i, parseErr)
    }
  prices[i] = price
	}
	return prices,nil
}