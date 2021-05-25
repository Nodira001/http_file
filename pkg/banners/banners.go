package banners

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"sync"
)

type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

var NextBannerID int64 = 0

func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.items, nil
}
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Println(id)
	for _, banner := range s.items {
		if id == banner.ID {
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}
func (s *Service) Save(ctx context.Context, item *Banner, image multipart.File, handler *multipart.FileHeader) (*Banner, error) {
	log.Println(item)
	if item.ID == 0 {
		NextBannerID++
		item.ID = NextBannerID
		if image != nil || item.Image == "" {
			item.Image = fmt.Sprint(NextBannerID) + path.Ext(handler.Filename)
			f, err := os.OpenFile("web/banners/"+fmt.Sprint(NextBannerID)+path.Ext(handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
			s.items = append(s.items, item)
			defer f.Close()
			io.Copy(f, image)
			if err != nil {
				log.Println(err)
			}
		}
		return item, nil
	}

	banner, err := s.ByID(ctx, item.ID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if banner != nil {
		banner.Title = item.Title
		banner.Link = item.Link
		banner.Content = item.Content
		banner.Button = item.Button

		if image != nil {
			banner.Image = fmt.Sprint(item.ID) + path.Ext(handler.Filename)
			f, err := os.OpenFile("web/banners/"+banner.Image, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			io.Copy(f, image)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return banner, nil
}
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	existent, err := s.ByID(ctx, id)
	if err != nil {
		return nil, err
	}
	b := []*Banner{}
	for _, banner := range s.items {
		if id != banner.ID {
			b = append(b, banner)
		}
	}
	if b == nil {
		b = []*Banner{}
	}
	s.items = b
	return existent, nil
}
