package command

import (
	"fmt"
	"github.com/aivyss/password-manager/repository"
	"github.com/urfave/cli/v2"
)

type AppVersionCommandHandler struct {
	appVersionRepository repository.AppVersionRepository
}

func (h *AppVersionCommandHandler) GetCurrentVersion(_ *cli.Context) error {
	version, err := h.appVersionRepository.GetLatestAppVersion()
	if err != nil {
		return err
	}

	fmt.Printf("[pwm] version: v%s\n", version)
	return nil
}

func (h *AppVersionCommandHandler) Command() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Action:  h.GetCurrentVersion,
	}
}

func NewAppVersionCommandHandler(appVersionRepository repository.AppVersionRepository) *AppVersionCommandHandler {
	return &AppVersionCommandHandler{
		appVersionRepository: appVersionRepository,
	}
}
