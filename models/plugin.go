package models

import (
	"net/http"
)

type PluginMetadata struct {
	Name        string
	Version     string
	Description string
}

// PluginConfig holds per-plugin configuration.
type PluginConfig struct {
	Enabled bool
	Options any
}

type PluginMiddleware struct {
	Auth          func() func(http.Handler) http.Handler
	OptionalAuth  func() func(http.Handler) http.Handler
	CorsAuth      func() func(http.Handler) http.Handler
	CSRF          func() func(http.Handler) http.Handler
	RateLimit     func() func(http.Handler) http.Handler
	EndpointHooks func() func(http.Handler) http.Handler
}

type PluginContext struct {
	Config     *Config
	Api        *Api
	EventBus   EventBus
	Middleware *PluginMiddleware
}

type PluginRouteMiddleware func(http.Handler) http.Handler

type PluginRouteHandler func() http.Handler

type PluginRoute struct {
	Method     string
	Path       string // Relative path, /auth is auto-prefixed
	Middleware []PluginRouteMiddleware
	Handler    PluginRouteHandler
}

type PluginRateLimit = RateLimitConfig

type BeforeCreateHook[T any] func(ctx *PluginContext, entity *T) error
type AfterCreateHook[T any] func(ctx *PluginContext, entity *T) error

type BeforeReadHook[T any] func(ctx *PluginContext) error
type AfterReadHook[T any] func(ctx *PluginContext, results *[]T) error

type BeforeUpdateHook[T any] func(ctx *PluginContext, existing *T, updatedData map[string]any) error
type AfterUpdateHook[T any] func(ctx *PluginContext, updated *T) error

type BeforeDeleteHook[T any] func(ctx *PluginContext, entity *T) error
type AfterDeleteHook[T any] func(ctx *PluginContext, entity *T) error

type PluginDatabaseHookOperations[T any] struct {
	BeforeCreate *BeforeCreateHook[T]
	AfterCreate  *AfterCreateHook[T]

	BeforeRead *BeforeReadHook[T]
	AfterRead  *AfterReadHook[T]

	BeforeUpdate *BeforeUpdateHook[T]
	AfterUpdate  *AfterUpdateHook[T]

	BeforeDelete *BeforeDeleteHook[T]
	AfterDelete  *AfterDeleteHook[T]
}

type PluginDatabaseHooks map[string]PluginDatabaseHookOperations[any]

type PluginEventHookPayload any

type PluginEventHookFunc func(ctx *PluginContext, payload PluginEventHookPayload) error

type PluginEventHooks map[string]PluginEventHookFunc

type TypedPluginEventHook[T any] func(ctx *PluginContext, payload T) error

type Plugin interface {
	Metadata() PluginMetadata
	SetMetadata(meta PluginMetadata)

	Config() PluginConfig
	SetConfig(cfg PluginConfig)

	Ctx() *PluginContext
	SetCtx(ctx *PluginContext)

	Init(ctx *PluginContext) error
	SetInit(fn func(ctx *PluginContext) error)

	Migrations() []any
	SetMigrations(migrations []any)

	Routes() []PluginRoute
	SetRoutes(routes []PluginRoute)

	RateLimit() *PluginRateLimit
	SetRateLimit(rateLimit *PluginRateLimit)

	DatabaseHooks() *PluginDatabaseHooks
	SetDatabaseHooks(hooks *PluginDatabaseHooks)

	EventHooks() *PluginEventHooks
	SetEventHooks(hooks *PluginEventHooks)

	Close() error
	SetClose(fn func() error)
}

type PluginOption func(p Plugin)

type BasePlugin struct {
	metadata      PluginMetadata
	config        PluginConfig
	ctx           *PluginContext
	init          func(ctx *PluginContext) error
	migrations    []any // Database migration structs (GORM models)
	routes        []PluginRoute
	rateLimit     *PluginRateLimit
	databaseHooks *PluginDatabaseHooks
	eventHooks    *PluginEventHooks
	close         func() error
}

func (p *BasePlugin) Metadata() PluginMetadata {
	return p.metadata
}

func (p *BasePlugin) SetMetadata(meta PluginMetadata) {
	p.metadata = meta
}

func (p *BasePlugin) Config() PluginConfig {
	return p.config
}

func (p *BasePlugin) SetConfig(cfg PluginConfig) {
	p.config = cfg
}

func (p *BasePlugin) Ctx() *PluginContext {
	return p.ctx
}

func (p *BasePlugin) SetCtx(ctx *PluginContext) {
	p.ctx = ctx
}

func (p *BasePlugin) Init(ctx *PluginContext) error {
	if p.init != nil {
		return p.init(ctx)
	}
	return nil
}

func (p *BasePlugin) SetInit(fn func(ctx *PluginContext) error) {
	p.init = fn
}

func (p *BasePlugin) Migrations() []any {
	return p.migrations
}

func (p *BasePlugin) SetMigrations(migrations []any) {
	p.migrations = migrations
}

func (p *BasePlugin) Routes() []PluginRoute {
	return p.routes
}

func (p *BasePlugin) SetRoutes(routes []PluginRoute) {
	p.routes = routes
}

func (p *BasePlugin) RateLimit() *PluginRateLimit {
	return p.rateLimit
}

func (p *BasePlugin) SetRateLimit(rateLimit *PluginRateLimit) {
	p.rateLimit = rateLimit
}

func (p *BasePlugin) DatabaseHooks() *PluginDatabaseHooks {
	return p.databaseHooks
}

func (p *BasePlugin) SetDatabaseHooks(hooks *PluginDatabaseHooks) {
	p.databaseHooks = hooks
}

func (p *BasePlugin) EventHooks() *PluginEventHooks {
	return p.eventHooks
}

func (p *BasePlugin) SetEventHooks(hooks *PluginEventHooks) {
	p.eventHooks = hooks
}

func (p *BasePlugin) Close() error {
	if p.close != nil {
		return p.close()
	}
	return nil
}

func (p *BasePlugin) SetClose(fn func() error) {
	p.close = fn
}
