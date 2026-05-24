package pagination

import (
    "math"
    "strconv"
    "strings"
)

type Pagination struct {
    Page     int    `json:"page" form:"page"`
    Limit    int    `json:"limit" form:"limit"`
    SortBy   string `json:"sort_by" form:"sort_by"`
    SortOrder string `json:"sort_order" form:"sort_order"`
    Search   string `json:"search" form:"search"`
}

type PageInfo struct {
    CurrentPage  int   `json:"current_page"`
    PerPage      int   `json:"per_page"`
    TotalItems   int64 `json:"total_items"`
    TotalPages   int   `json:"total_pages"`
    HasNext      bool  `json:"has_next"`
    HasPrevious  bool  `json:"has_previous"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    PageInfo   PageInfo    `json:"page_info"`
}

func NewPagination(page, limit int, sortBy, sortOrder, search string) *Pagination {
    if page < 1 {
        page = 1
    }
    
    if limit < 1 {
        limit = 10
    }
    
    if limit > 100 {
        limit = 100
    }
    
    if sortBy == "" {
        sortBy = "created_at"
    }
    
    sortOrder = strings.ToLower(sortOrder)
    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = "desc"
    }
    
    return &Pagination{
        Page:      page,
        Limit:     limit,
        SortBy:    sortBy,
        SortOrder: sortOrder,
        Search:    strings.TrimSpace(search),
    }
}

func (p *Pagination) GetOffset() int {
    return (p.Page - 1) * p.Limit
}

func (p *Pagination) GetLimit() int {
    return p.Limit
}

func (p *Pagination) GetSortClause() string {
    return p.SortBy + " " + p.SortOrder
}

func (p *Pagination) GetSearchPattern() string {
    if p.Search == "" {
        return ""
    }
    return "%" + p.Search + "%"
}

func (p *Pagination) HasSearch() bool {
    return p.Search != ""
}

func NewPageInfo(totalItems int64, page, limit int) PageInfo {
    totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
    
    return PageInfo{
        CurrentPage: page,
        PerPage:     limit,
        TotalItems:  totalItems,
        TotalPages:  totalPages,
        HasNext:     page < totalPages,
        HasPrevious: page > 1,
    }
}

func NewPaginatedResponse(data interface{}, totalItems int64, page, limit int) PaginatedResponse {
    return PaginatedResponse{
        Data:     data,
        PageInfo: NewPageInfo(totalItems, page, limit),
    }
}

func ParsePaginationFromQuery(queryParams map[string][]string) *Pagination {
    page := 1
    limit := 10
    sortBy := "created_at"
    sortOrder := "desc"
    search := ""
    
    if pageStr, ok := queryParams["page"]; ok && len(pageStr) > 0 {
        if p, err := strconv.Atoi(pageStr[0]); err == nil && p > 0 {
            page = p
        }
    }
    
    if limitStr, ok := queryParams["limit"]; ok && len(limitStr) > 0 {
        if l, err := strconv.Atoi(limitStr[0]); err == nil && l > 0 && l <= 100 {
            limit = l
        }
    }
    
    if sortByStr, ok := queryParams["sort_by"]; ok && len(sortByStr) > 0 {
        sortBy = sortByStr[0]
    }
    
    if sortOrderStr, ok := queryParams["sort_order"]; ok && len(sortOrderStr) > 0 {
        order := strings.ToLower(sortOrderStr[0])
        if order == "asc" || order == "desc" {
            sortOrder = order
        }
    }
    
    if searchStr, ok := queryParams["search"]; ok && len(searchStr) > 0 {
        search = searchStr[0]
    }
    
    return NewPagination(page, limit, sortBy, sortOrder, search)
}

type CursorPagination struct {
    Cursor string `json:"cursor" form:"cursor"`
    Limit  int    `json:"limit" form:"limit"`
}

func NewCursorPagination(cursor string, limit int) *CursorPagination {
    if limit < 1 {
        limit = 10
    }
    if limit > 100 {
        limit = 100
    }
    
    return &CursorPagination{
        Cursor: cursor,
        Limit:  limit,
    }
}

func (c *CursorPagination) GetLimit() int {
    return c.Limit
}

func (c *CursorPagination) HasCursor() bool {
    return c.Cursor != ""
}

type PaginationConfig struct {
    DefaultLimit   int
    MaxLimit       int
    DefaultSortBy  string
    DefaultSortOrder string
}

func DefaultPaginationConfig() PaginationConfig {
    return PaginationConfig{
        DefaultLimit:   10,
        MaxLimit:       100,
        DefaultSortBy:  "created_at",
        DefaultSortOrder: "desc",
    }
}

func NewPaginationWithConfig(page, limit int, sortBy, sortOrder, search string, config PaginationConfig) *Pagination {
    if page < 1 {
        page = 1
    }
    
    if limit < 1 {
        limit = config.DefaultLimit
    }
    
    if limit > config.MaxLimit {
        limit = config.MaxLimit
    }
    
    if sortBy == "" {
        sortBy = config.DefaultSortBy
    }
    
    sortOrder = strings.ToLower(sortOrder)
    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = config.DefaultSortOrder
    }
    
    return &Pagination{
        Page:      page,
        Limit:     limit,
        SortBy:    sortBy,
        SortOrder: sortOrder,
        Search:    strings.TrimSpace(search),
    }
}

// Helper function to build pagination SQL query
type PaginatedQuery struct {
    Query     string
    CountQuery string
    Args      []interface{}
    Offset    int
    Limit     int
}

func BuildPaginatedQuery(baseQuery string, countQuery string, pagination *Pagination, args ...interface{}) PaginatedQuery {
    offset := pagination.GetOffset()
    limit := pagination.GetLimit()
    
    if !strings.Contains(strings.ToUpper(baseQuery), "ORDER BY") {
        baseQuery += " ORDER BY " + pagination.GetSortClause()
    }
    
    baseQuery += " LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
    
    allArgs := append(args, limit, offset)
    
    return PaginatedQuery{
        Query:      baseQuery,
        CountQuery: countQuery,
        Args:       allArgs,
        Offset:     offset,
        Limit:      limit,
    }
}

type PaginationExample struct {
    Name string
    Description string
}

func GetPaginationExample() map[string]interface{} {
    return map[string]interface{}{
        "query_parameters": map[string]string{
            "page":         "Page number (default: 1)",
            "limit":        "Items per page (default: 10, max: 100)",
            "sort_by":      "Field to sort by (default: created_at)",
            "sort_order":   "Sort order: asc or desc (default: desc)",
            "search":       "Search term (optional)",
        },
        "example_request": "/api/v1/users?page=2&limit=20&sort_by=name&sort_order=asc&search=john",
        "example_response": map[string]interface{}{
            "data": []interface{}{},
            "page_info": map[string]interface{}{
                "current_page":   2,
                "per_page":       20,
                "total_items":    150,
                "total_pages":    8,
                "has_next":       true,
                "has_previous":   true,
            },
        },
    }
}