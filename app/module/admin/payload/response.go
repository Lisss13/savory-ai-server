package payload

import "time"

// StatsResp - статистика системы
type StatsResp struct {
	TotalUsers         int64            `json:"totalUsers"`
	ActiveUsers        int64            `json:"activeUsers"`
	TotalOrganizations int64            `json:"totalOrganizations"`
	TotalRestaurants   int64            `json:"totalRestaurants"`
	TotalDishes        int64            `json:"totalDishes"`
	TotalTables        int64            `json:"totalTables"`
	TotalQuestions     int64            `json:"totalQuestions"`
	ActiveSubscriptions int64           `json:"activeSubscriptions"`
	RecentActivity     []ActivityResp   `json:"recentActivity"`
}

// ActivityResp - недавняя активность
type ActivityResp struct {
	ID         uint      `json:"id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entityType"`
	EntityID   uint      `json:"entityId"`
	AdminName  string    `json:"adminName"`
	CreatedAt  time.Time `json:"createdAt"`
}

// AdminUserResp - пользователь для админ-панели
type AdminUserResp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Company   string    `json:"company"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

// AdminUsersResp - список пользователей
type AdminUsersResp struct {
	Users      []AdminUserResp `json:"users"`
	TotalCount int64           `json:"totalCount"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
}

// AdminOrganizationResp - организация для админ-панели
type AdminOrganizationResp struct {
	ID                  uint      `json:"id"`
	Name                string    `json:"name"`
	Phone               string    `json:"phone"`
	AdminID             uint      `json:"adminId"`
	AdminName           string    `json:"adminName"`
	AdminEmail          string    `json:"adminEmail"`
	RestaurantsCount    int64     `json:"restaurantsCount"`
	UsersCount          int64     `json:"usersCount"`
	HasActiveSubscription bool    `json:"hasActiveSubscription"`
	CreatedAt           time.Time `json:"createdAt"`
}

// AdminOrganizationsResp - список организаций
type AdminOrganizationsResp struct {
	Organizations []AdminOrganizationResp `json:"organizations"`
	TotalCount    int64                   `json:"totalCount"`
	Page          int                     `json:"page"`
	PageSize      int                     `json:"pageSize"`
}

// AdminDishResp - блюдо для модерации
type AdminDishResp struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	Image          string    `json:"image"`
	RestaurantID   uint      `json:"restaurantId"`
	RestaurantName string    `json:"restaurantName"`
	CreatedAt      time.Time `json:"createdAt"`
}

// AdminDishesResp - список блюд для модерации
type AdminDishesResp struct {
	Dishes     []AdminDishResp `json:"dishes"`
	TotalCount int64           `json:"totalCount"`
	Page       int             `json:"page"`
	PageSize   int             `json:"pageSize"`
}

// AdminLogResp - лог действия администратора
type AdminLogResp struct {
	ID         uint      `json:"id"`
	AdminID    uint      `json:"adminId"`
	AdminName  string    `json:"adminName"`
	AdminEmail string    `json:"adminEmail"`
	Action     string    `json:"action"`
	EntityType string    `json:"entityType"`
	EntityID   uint      `json:"entityId"`
	Details    string    `json:"details"`
	IPAddress  string    `json:"ipAddress"`
	CreatedAt  time.Time `json:"createdAt"`
}

// AdminLogsResp - список логов
type AdminLogsResp struct {
	Logs       []AdminLogResp `json:"logs"`
	TotalCount int64          `json:"totalCount"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
}
