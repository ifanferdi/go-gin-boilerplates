package category

import (
	"context"
	"time"

	"go-gin-boilerplate/internal/helpers"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type CategoryGRPCServer interface {
	FindAll(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	FindByID(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	Create(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	Update(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
	Delete(ctx context.Context, req *structpb.Struct) (*structpb.Struct, error)
}

type categoryGRPCServer struct {
	service CategoryService
}

func NewCategoryGRPCServer(service CategoryService) CategoryGRPCServer {
	return &categoryGRPCServer{service: service}
}

func (g *categoryGRPCServer) FindAll(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	fields := req.GetFields()
	params := FindAllCategoryDTO{
		Page:      int(fields["page"].GetNumberValue()),
		Limit:     int(fields["limit"].GetNumberValue()),
		Search:    fields["q"].GetStringValue(),
		Name:      fields["name"].GetStringValue(),
		OrderBy:   fields["orderBy"].GetStringValue(),
		ProductID: uint(fields["productId"].GetNumberValue()),
	}

	result, err := g.service.FindAll(params)
	if err != nil {
		return nil, toGRPCError(err)
	}

	items := make([]interface{}, 0, len(result.Data))
	for _, c := range result.Data {
		items = append(items, categoryToMap(&c))
	}

	return structpb.NewStruct(map[string]interface{}{
		"page":       result.Page,
		"limit":      result.Limit,
		"total":      result.Total,
		"totalPages": result.TotalPages,
		"data":       items,
	})
}

func (g *categoryGRPCServer) FindByID(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	id := uint(req.GetFields()["id"].GetNumberValue())
	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Column id: ID cannot be empty.")
	}

	item, err := g.service.FindByID(id)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(categoryToMap(item))
}

func (g *categoryGRPCServer) Create(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	fields := req.GetFields()
	name := fields["name"].GetStringValue()
	var desc *string
	if v, ok := fields["description"]; ok && v.GetStringValue() != "" {
		s := v.GetStringValue()
		desc = &s
	}

	item, err := g.service.Create(CreateCategoryDTO{
		Name:        name,
		Description: desc,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(categoryToMap(item))
}

func (g *categoryGRPCServer) Update(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	fields := req.GetFields()
	id := uint(fields["id"].GetNumberValue())
	name := fields["name"].GetStringValue()
	var desc *string
	if v, ok := fields["description"]; ok {
		s := v.GetStringValue()
		desc = &s
	}

	item, err := g.service.Update(UpdateCategoryDTO{
		ID:          id,
		Name:        name,
		Description: desc,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(categoryToMap(item))
}

func (g *categoryGRPCServer) Delete(_ context.Context, req *structpb.Struct) (*structpb.Struct, error) {
	id := uint(req.GetFields()["id"].GetNumberValue())
	if err := g.service.Delete(id); err != nil {
		return nil, toGRPCError(err)
	}

	return structpb.NewStruct(map[string]interface{}{
		"message": "Success.",
	})
}

func categoryToMap(c *Category) map[string]interface{} {
	var desc interface{}
	if c.Description != nil {
		desc = *c.Description
	}
	return map[string]interface{}{
		"id":          c.ID,
		"name":        c.Name,
		"description": desc,
		"createdAt":   c.CreatedAt.Format(time.RFC3339),
		"updatedAt":   c.UpdatedAt.Format(time.RFC3339),
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

var CategoryServiceDesc = grpc.ServiceDesc{
	ServiceName: "category.v1.CategoryService",
	HandlerType: (*CategoryGRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAll",
			Handler:    grpcHandler(func(srv CategoryGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.FindAll(ctx, req) }),
		},
		{
			MethodName: "FindByID",
			Handler:    grpcHandler(func(srv CategoryGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.FindByID(ctx, req) }),
		},
		{
			MethodName: "Create",
			Handler:    grpcHandler(func(srv CategoryGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.Create(ctx, req) }),
		},
		{
			MethodName: "Update",
			Handler:    grpcHandler(func(srv CategoryGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.Update(ctx, req) }),
		},
		{
			MethodName: "Delete",
			Handler:    grpcHandler(func(srv CategoryGRPCServer, ctx context.Context, req *structpb.Struct) (*structpb.Struct, error) { return srv.Delete(ctx, req) }),
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "docs/proto/category.proto",
}

func grpcHandler(fn func(CategoryGRPCServer, context.Context, *structpb.Struct) (*structpb.Struct, error)) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		in := new(structpb.Struct)
		if err := dec(in); err != nil {
			return nil, err
		}
		if interceptor == nil {
			return fn(srv.(CategoryGRPCServer), ctx, in)
		}
		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: "/category.v1.CategoryService",
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return fn(srv.(CategoryGRPCServer), ctx, req.(*structpb.Struct))
		}
		return interceptor(ctx, in, info, handler)
	}
}
