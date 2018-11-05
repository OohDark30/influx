package http_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/EMCECS/influx"
	"github.com/EMCECS/influx/http"
	"github.com/EMCECS/influx/inmem"
	"github.com/EMCECS/influx/mock"
	_ "github.com/EMCECS/influx/query/builtin"
	"github.com/EMCECS/influx/task"
	"github.com/EMCECS/influx/task/backend"
	tmock "github.com/EMCECS/influx/task/mock"
	"github.com/EMCECS/influx/task/servicetest"
	"go.uber.org/zap/zaptest"
)

func httpTaskServiceFactory(t *testing.T) (*servicetest.System, context.CancelFunc) {
	store := backend.NewInMemStore()
	rrw := backend.NewInMemRunReaderWriter()
	sch := tmock.NewScheduler()

	ctx, cancel := context.WithCancel(context.Background())

	backingTS := task.PlatformAdapter(store, rrw, sch)

	i := inmem.NewService()

	mappingService := mock.NewUserResourceMappingService()
	h := http.NewAuthenticationHandler()
	h.AuthorizationService = i
	th := http.NewTaskHandler(mappingService, zaptest.NewLogger(t))
	th.TaskService = backingTS
	th.AuthorizationService = i
	h.Handler = th

	org := &platform.Organization{Name: t.Name() + "_org"}
	if err := i.CreateOrganization(ctx, org); err != nil {
		t.Fatal(err)
	}
	user := &platform.User{Name: t.Name() + "_user"}
	if err := i.CreateUser(ctx, user); err != nil {
		t.Fatal(err)
	}
	auth := platform.Authorization{UserID: user.ID}
	if err := i.CreateAuthorization(ctx, &auth); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(h)
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	tsFunc := func() platform.TaskService {
		return http.TaskService{
			Addr:  server.URL,
			Token: auth.Token,
		}
	}

	cFunc := func() (o, u platform.ID, tok string, err error) {
		return org.ID, user.ID, auth.Token, nil
	}

	return &servicetest.System{
		S:               store,
		LR:              rrw,
		LW:              rrw,
		Ctx:             ctx,
		TaskServiceFunc: tsFunc,
		CredsFunc:       cFunc,
	}, cancel
}

func TestTaskService(t *testing.T) {
	t.Parallel()

	servicetest.TestTaskService(t, httpTaskServiceFactory)
}
