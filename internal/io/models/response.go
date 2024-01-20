package models

import "em_test/internal/db/repo"

type People struct {
  Data       []repo.Person `json:"data"`
  TotalCount int           `json:"total_count"`
  Limit      int           `json:"limit"`
  Page       int           `json:"page"`
  PrevPage   int           `json:"prev_page,omitempty"`
  NextPage   int           `json:"next_page,omitempty"`
}

type ErrResp struct {
  Err string
}

type OkResp struct {
  Msg string
}
