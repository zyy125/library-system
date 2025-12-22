package service

import (
	"context"
	"errors"
	"library-system/common"
	"library-system/dto/request"
	"library-system/dto/response"
	"library-system/model"
	"library-system/repository"

	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo *repository. CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetCategoryList(ctx context.Context, req *request.GetCategoryListRequest) (*response.GetCategoryListResponse, error) {
	if req.IncludeCount {
		categories, counts, err := s. categoryRepo.GetCategoryListWithBookCount(ctx)
		if err != nil {
			return nil, err
		}

		items := make([]response.CategoryItem, len(categories))
		for i, cat := range categories {
			items[i] = response.CategoryItem{
				ID:          cat.ID,
				Name:        cat.Name,
				Description: cat.Description,
				BookCount:   counts[i],
				ParentID:    cat.ParentID,
				CreatedAt:   cat.CreatedAt,
			}
		}

		return &response.GetCategoryListResponse{Categories: items}, nil
	}

	categories, err := s. categoryRepo.GetCategoryList(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]response. CategoryItem, len(categories))
	for i, cat := range categories {
		items[i] = response.CategoryItem{
			ID:          cat.ID,
			Name:        cat. Name,
			Description: cat.Description,
			ParentID:    cat.ParentID,
			CreatedAt:   cat.CreatedAt,
		}
	}

	return &response. GetCategoryListResponse{Categories: items}, nil
}

func (s *CategoryService) GetCategoryDetail(ctx context.Context, id uint) (*response.GetCategoryDetailResponse, error) {
	category, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrCategoryNotFound
		}
		return nil, err
	}

	bookCount, err := s.categoryRepo.GetBookCountByCategoryID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取子分类
	children, err := s.categoryRepo.GetChildCategories(ctx, id)
	if err != nil {
		return nil, err
	}

	childItems := make([]response.CategoryChild, len(children))
	for i, child := range children {
		count, _ := s.categoryRepo.GetBookCountByCategoryID(ctx, child.ID)
		childItems[i] = response. CategoryChild{
			ID:        child. ID,
			Name:      child.Name,
			BookCount: count,
		}
	}

	resp := &response.GetCategoryDetailResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		BookCount:   bookCount,
		ParentID:    category.ParentID,
		CreatedAt:   category.CreatedAt,
	}

	if len(childItems) > 0 {
		resp.Children = childItems
	}

	return resp, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *request.CreateCategoryRequest) (*response.CreateCategoryResponse, error) {
	// 检查名称是否已存在
	if _, err := s.categoryRepo.GetCategoryByName(ctx, req.Name); err == nil {
		return nil, &common.BizError{
			Code:       409,
			Message:     "分类名称已存在",
			HTTPStatus: 409,
		}
	} else if !errors. Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 如果有父分类，检查父分类是否存在
	if req.ParentID != nil {
		if _, err := s. categoryRepo.GetCategoryByID(ctx, *req.ParentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, &common. BizError{
					Code:       400,
					Message:     "父分类不存在",
					HTTPStatus: 400,
				}
			}
			return nil, err
		}
	}

	category := model.Category{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	if req.Description != nil {
		category.Description = *req.Description
	}

	if err := s. categoryRepo.CreateCategory(ctx, &category); err != nil {
		return nil, err
	}

	return &response. CreateCategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category. ParentID,
		CreatedAt:    category.CreatedAt,
	}, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id uint, req *request.UpdateCategoryRequest) (*response.UpdateCategoryResponse, error) {
	category, err := s. categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		if errors. Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ErrCategoryNotFound
		}
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Name != nil {
		// 检查新名称是否与其他分类冲突
		existing, err := s. categoryRepo.GetCategoryByName(ctx, *req.Name)
		if err == nil && existing.ID != id {
			return nil, &common.BizError{
				Code:       409,
				Message:     "分类名称已存在",
				HTTPStatus: 409,
			}
		}
		updates["name"] = *req.Name
	}

	if req. Description != nil {
		updates["description"] = *req.Description
	}

	if req.ParentID != nil {
		// 检查父分类是否存在
		if _, err := s. categoryRepo.GetCategoryByID(ctx, *req.ParentID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, &common. BizError{
					Code:       400,
					Message:     "父分类不存在",
					HTTPStatus: 400,
				}
			}
			return nil, err
		}
		// 不能将自己设为父分类
		if *req.ParentID == id {
			return nil, &common.BizError{
				Code:       400,
				Message:    "不能将自己设为父分类",
				HTTPStatus: 400,
			}
		}
		updates["parent_id"] = *req. ParentID
	}

	if len(updates) == 0 {
		return nil, &common.BizError{
			Code:        400,
			Message:    "没有需要更新的字段",
			HTTPStatus: 400,
		}
	}

	if err := s.categoryRepo.UpdateCategory(ctx, id, updates); err != nil {
		return nil, err
	}

	// 重新获取更新后的分类
	updatedCategory, _ := s.categoryRepo.GetCategoryByID(ctx, id)

	return &response. UpdateCategoryResponse{
		ID:        category.ID,
		Name:      &updatedCategory.Name,
		UpdatedAt: updatedCategory.UpdatedAt,
	}, nil
}

func (s *CategoryService) DeleteCategory(ctx context. Context, id uint) error {
	// 检查分类是否存在
	if _, err := s.categoryRepo.GetCategoryByID(ctx, id); err != nil {
		if errors. Is(err, gorm.ErrRecordNotFound) {
			return common. ErrCategoryNotFound
		}
		return err
	}

	// 检查是否有子分类
	hasChildren, err := s.categoryRepo.HasChildren(ctx, id)
	if err != nil {
		return err
	}
	if hasChildren {
		return (&common.BizError{
			Code:       400,
			Message:    "无法删除该分类",
			HTTPStatus: 400,
		}).WithDetails(map[string]interface{}{
			"reason": "分类下存在子分类",
		})
	}

	// 检查分类下是否有图书
	hasBooks, bookCount, err := s.categoryRepo.HasBooks(ctx, id)
	if err != nil {
		return err
	}
	if hasBooks {
		return (&common.BizError{
			Code:       400,
			Message:    "无法删除该分类",
			HTTPStatus: 400,
		}).WithDetails(map[string]interface{}{
			"reason":      "分类下存在图书",
			"book_count": bookCount,
		})
	}

	return s.categoryRepo. DeleteCategory(ctx, id)
}