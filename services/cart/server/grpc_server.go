package server

import (
	"context"

	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/grpc_utils"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/pkg/tools/logger"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/pkg/models"
	proto "github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/proto/cart"
	"github.com/go-park-mail-ru/2021_1_DuckLuck/services/cart/server/errors"

	"github.com/golang/protobuf/ptypes/empty"
)

type CartServer struct {
	CartUCase cart.UseCase
}

func NewCartServer(cartUCase cart.UseCase) *CartServer {
	return &CartServer{
		CartUCase: cartUCase,
	}
}

func (s *CartServer) AddProduct(ctx context.Context,
	cartArticle *proto.ReqCartArticle) (*empty.Empty, error) {
	var err error
	defer func() {
		requireId := grpc_utils.MustGetRequireId(ctx)
		if err != nil {
			logger.LogError("cart_service_handler", "AddProduct", requireId, err)
		}
	}()

	err = s.CartUCase.AddProduct(cartArticle.UserId, &models.CartArticle{
		ProductPosition: models.ProductPosition{
			Count: cartArticle.Position.Count,
		},
		ProductIdentifier: models.ProductIdentifier{
			ProductId: cartArticle.ProductId,
		},
	})

	if err != nil {
		return &empty.Empty{}, errors.CreateError(err)
	}

	return &empty.Empty{}, nil
}

func (s *CartServer) DeleteProduct(ctx context.Context,
	productIdentifier *proto.ReqProductIdentifier) (*empty.Empty, error) {
	var err error
	defer func() {
		requireId := grpc_utils.MustGetRequireId(ctx)
		if err != nil {
			logger.LogError("cart_service_handler", "DeleteProduct", requireId, err)
		}
	}()

	err = s.CartUCase.DeleteProduct(productIdentifier.UserId,
		&models.ProductIdentifier{
			ProductId: productIdentifier.ProductId,
		})

	if err != nil {
		return &empty.Empty{}, errors.CreateError(err)
	}

	return &empty.Empty{}, nil
}

func (s *CartServer) ChangeProduct(ctx context.Context,
	cartArticle *proto.ReqCartArticle) (*empty.Empty, error) {
	var err error
	defer func() {
		requireId := grpc_utils.MustGetRequireId(ctx)
		if err != nil {
			logger.LogError("cart_service_handler", "ChangeProduct", requireId, err)
		}
	}()

	err = s.CartUCase.ChangeProduct(cartArticle.UserId, &models.CartArticle{
		ProductPosition: models.ProductPosition{
			Count: cartArticle.Position.Count,
		},
		ProductIdentifier: models.ProductIdentifier{
			ProductId: cartArticle.ProductId,
		},
	})

	if err != nil {
		return &empty.Empty{}, errors.CreateError(err)
	}

	return &empty.Empty{}, nil
}

func (s *CartServer) GetPreviewCart(ctx context.Context,
	userId *proto.ReqUserId) (*proto.Cart, error) {
	var err error
	defer func() {
		requireId := grpc_utils.MustGetRequireId(ctx)
		if err != nil {
			logger.LogError("cart_service_handler", "GetPreviewCart", requireId, err)
		}
	}()

	userCart, err := s.CartUCase.GetPreviewCart(userId.UserId)

	if err != nil {
		return nil, errors.CreateError(err)
	}

	userCartProto := &proto.Cart{}
	userCartProto.Products = make(map[uint64]*proto.ProductPosition)
	for key, item := range userCart.Products {
		userCartProto.Products[key] = &proto.ProductPosition{
			Count: item.Count,
		}
	}

	return userCartProto, nil
}

func (s *CartServer) DeleteCart(ctx context.Context,
	userId *proto.ReqUserId) (*empty.Empty, error) {
	var err error
	defer func() {
		requireId := grpc_utils.MustGetRequireId(ctx)
		if err != nil {
			logger.LogError("cart_service_handler", "DeleteCart", requireId, err)
		}
	}()

	err = s.CartUCase.DeleteCart(userId.UserId)

	if err != nil {
		return &empty.Empty{}, errors.CreateError(err)
	}

	return &empty.Empty{}, nil
}
