package product

import (
	"context"
	"time"

	"go-gin-boilerplate/internal/helpers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type ProductGRPCServer interface {
	FindAll(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	FindByID(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	Create(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	Update(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	Delete(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
}

type productGRPCServer struct {
	service ProductService
}

func NewProductGRPCServer(service ProductService) ProductGRPCServer {
	return &productGRPCServer{service: service}
}

func (g *productGRPCServer) FindAll(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	fields := req.GetFields()
	params := FindAllProductDTO{
		Page:       int(fields["page"].GetNumberValue()),
		Limit:      int(fields["limit"].GetNumberValue()),
		Search:     fields["q"].GetStringValue(),
		SKU:        fields["sku"].GetStringValue(),
		Barcode:    fields["barcode"].GetStringValue(),
		CategoryID: uint(fields["categoryId"].GetNumberValue()),
		OrderBy:    fields["orderBy"].GetStringValue(),
	}

	result, err := g.service.FindAll(params)
	if err != nil {
		return nil, toGRPCError(err)
	}

	items := make([]interface{}, 0, len(result.Data))
	for _, p := range result.Data {
		items = append(items, productToMap(&p))
	}

	return structpb.NewStruct(map[string]interface{}{
		"page":       result.Page,
		"limit":      result.Limit,
		"total":      result.Total,
		"totalPages": result.TotalPages,
		"data":       items,
	})
}

func (g *productGRPCServer) FindByID(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	id := uint(req.GetFields()["id"].GetNumberValue())
	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Column id: ID cannot be empty.")
	}

	item, err := g.service.FindByID(id)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(productToMap(item))
}

func (g *productGRPCServer) Create(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	fields := req.GetFields()
	dto := CreateProductDTO{
		Name:        fields["name"].GetStringValue(),
		Price:       int(fields["price"].GetNumberValue()),
		Stock:       int(fields["stock"].GetNumberValue()),
		CategoryIDs: extractUintList(fields["categoryIds"]),
	}
	if v, ok := fields["description"]; ok && v.GetStringValue() != "" {
		s := v.GetStringValue()
		dto.Description = &s
	}
	if v, ok := fields["sku"]; ok && v.GetStringValue() != "" {
		s := v.GetStringValue()
		dto.SKU = &s
	}
	if v, ok := fields["barcode"]; ok && v.GetStringValue() != "" {
		s := v.GetStringValue()
		dto.Barcode = &s
	}

	item, err := g.service.Create(dto)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(productToMap(item))
}

func (g *productGRPCServer) Update(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	fields := req.GetFields()
	dto := UpdateProductDTO{
		ID:    uint(fields["id"].GetNumberValue()),
		Name:  fields["name"].GetStringValue(),
		Price: int(fields["price"].GetNumberValue()),
		Stock: int(fields["stock"].GetNumberValue()),
	}
	if v, ok := fields["categoryIds"]; ok {
		ids := extractUintList(v)
		dto.CategoryIDs = &ids
	}

	item, err := g.service.Update(dto)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(productToMap(item))
}

func (g *productGRPCServer) Delete(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	id := uint(req.GetFields()["id"].GetNumberValue())
	if err := g.service.Delete(id); err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(map[string]interface{}{
		"message": "Success.",
	})
}

func extractUintList(val *structpb.Value) []uint {
	if val == nil || val.GetListValue() == nil {
		return nil
	}
	vals := val.GetListValue().GetValues()
	out := make([]uint, 0, len(vals))
	for _, v := range vals {
		if n := uint(v.GetNumberValue()); n > 0 {
			out = append(out, n)
		}
	}
	return out
}

func productToMap(p *Product) map[string]interface{} {
	cats := make([]interface{}, 0, len(p.Categories))
	for _, c := range p.Categories {
		cats = append(cats, map[string]interface{}{
			"id":   c.ID,
			"name": c.Name,
		})
	}
	return map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"price":      p.Price,
		"stock":      p.Stock,
		"isActive":   p.IsActive,
		"categories": cats,
		"createdAt":  p.CreatedAt.Format(time.RFC3339),
		"updatedAt":  p.UpdatedAt.Format(time.RFC3339),
	}
}

func toGRPCError(err error) error {
	if appErr, ok := helpers.IsAppError(err); ok {
		switch appErr.StatusCode {
		case 400:
			return status.Error(codes.InvalidArgument, appErr.Message)
		case 404:
			return status.Error(codes.NotFound, appErr.Message)
		case 409:
			return status.Error(codes.AlreadyExists, appErr.Message)
		}
	}
	return status.Error(codes.Internal, err.Error())
}

var ProductServiceDesc = grpc.ServiceDesc{
	ServiceName: "product.v1.ProductService",
	HandlerType: (*ProductGRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAll",
			Handler:    grpcHandler(func(srv ProductGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.FindAll(ctx, req) }),
		},
		{
			MethodName: "FindByID",
			Handler:    grpcHandler(func(srv ProductGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.FindByID(ctx, req) }),
		},
		{
			MethodName: "Create",
			Handler:    grpcHandler(func(srv ProductGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.Create(ctx, req) }),
		},
		{
			MethodName: "Update",
			Handler:    grpcHandler(func(srv ProductGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.Update(ctx, req) }),
		},
		{
			MethodName: "Delete",
			Handler:    grpcHandler(func(srv ProductGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.Delete(ctx, req) }),
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "docs/proto/product.proto",
}

func grpcHandler(fn func(ProductGRPCServer, context.Context, *structpb.Struct) (*structpb.Struct, error)) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := new(structpb.Struct)
		if err := dec(in); err != nil {
			return nil, err
		}
		if interceptor == nil {
			return fn(srv.(ProductGRPCServer), ctx, in)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: "/product.v1.ProductService",
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return fn(srv.(ProductGRPCServer), ctx, req.(*structpb.Struct))
		}
		return interceptor(ctx, in, info, handler)
	}
}
