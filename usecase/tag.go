package usecase

import (
	"errors"
	"fmt"

	"github.com/masibw/blog-server/domain/dto"
	"github.com/masibw/blog-server/domain/entity"
	"github.com/masibw/blog-server/domain/repository"
)

type TagUseCase struct {
	tagRepository repository.Tag
}

func NewTagUseCase(tagRepository repository.Tag) *TagUseCase {
	return &TagUseCase{tagRepository: tagRepository}
}

func (p *TagUseCase) StoreTag(tagDTO *dto.TagDTO) (*dto.TagDTO, error) {
	var tag *entity.Tag
	var err error

	tag, err = p.tagRepository.FindByName(tagDTO.Name)
	if err != nil && !errors.Is(err, entity.ErrTagNotFound) {
		return nil, fmt.Errorf("store tag name=%v: %w", tagDTO.Name, err)
	}
	if tag != nil {
		return nil, fmt.Errorf("store tag name=%v: %w", tagDTO.Name, entity.ErrTagNameAlreadyExisted)
	}

	tag = entity.NewTag(tagDTO.Name)

	err = p.tagRepository.Store(tag)
	if err != nil {
		return nil, fmt.Errorf("store tag name=%v: %w", tagDTO.Name, err)
	}

	return tag.ConvertToDTO(), nil
}

func (p *TagUseCase) GetTags(offset, pageSize int, condition string, params []interface{}) (tagDTOs []*dto.TagDTO, count int, err error) {
	var tags []*entity.Tag
	tags, err = p.tagRepository.FindAll(offset, pageSize, condition, params)
	if err != nil {
		err = fmt.Errorf("get tags: %w", err)
		return
	}
	count, err = p.tagRepository.Count()
	if err != nil {
		err = fmt.Errorf("count tags: %w", err)
		return
	}
	for _, tag := range tags {
		tagDTOs = append(tagDTOs, tag.ConvertToDTO())
	}

	return
}

func (p *TagUseCase) GetTag(id string) (tagDTO *dto.TagDTO, err error) {
	var tag *entity.Tag
	tag, err = p.tagRepository.FindByID(id)
	if err != nil {
		err = fmt.Errorf("get tag: %w", err)
		return
	}
	tagDTO = tag.ConvertToDTO()
	return
}

func (p *TagUseCase) DeleteTag(id string) (err error) {
	err = p.tagRepository.Delete(id)
	if err != nil {
		err = fmt.Errorf("delete tag: %w", err)
		return
	}
	return nil
}
