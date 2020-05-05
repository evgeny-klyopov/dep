package config

type Config struct {
	Repository         *string            `json:"repository"`
	LocalObjectPath    *[]string          `json:"local_object_path"`
	Hosts              []Host             `json:"hosts" validate:"required,dive,required"`
	KeepReleases       *int8              `json:"keep_releases"`

	OrderTasks		orderTasks	`json:"order_tasks" validate:"required,dive,required"`
	//DeployTasksOrder   []string           `json:"deploy_tasks_order" validate:"required,min=1"`
	//RollbackTasksOrder *[]string           `json:"rollback_tasks_order"`

	Writable           *[]string          `json:"writable"`
	Variables          *map[string]string `json:"variables"`
	Shared             *[]shared          `json:"shared"`
	Tasks              *tasks             `json:"tasks" validate:""`
	Notifications      *notifications     `json:"notifications" validate:""`
}

type orderTasks struct {
	Deploy []string           `json:"deploy" validate:"required,min=1"`
	Rollback  *[]string           `json:"rollback"`
}

type Host struct {
	Host       string  `json:"host" validate:"required,min=5"`
	Branch     *string `json:"branch"`
	Stage      string  `json:"stage" validate:"required,min=1"`
	User       string  `json:"user" validate:"required,min=1"`
	Port       int     `json:"port" validate:"required,min=1"`
	DeployPath string  `json:"deploy_path" validate:"required,min=1"`
}
type shared struct {
	Path  string `json:"path" validate:"required,min=1"`
	IsDir bool   `json:"is_dir" validate:"required"`
}

type task struct {
	Name    string `json:"name" validate:"required,min=1"`
	Command string `json:"command" validate:"required,min=1"`
}

type tasks struct {
	Remote *[]task `json:"remote" validate:"required,dive,required"`
	Local  *[]task `json:"local" validate:"required,dive,required"`
}
type telegram struct {
	UseProxy bool   `json:"use_proxy"`
	Proxy    string `json:"proxy"`
	ChatId   int64  `json:"chat_id" validate:"required,min=1"`
	Token    string `json:"token" validate:"required,min=1"`
}
type notifications struct {
	Telegram *[]telegram `json:"telegram" validate:"required,dive,required"`
}
