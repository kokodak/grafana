package legacy_storage

import (
	"errors"
	"fmt"
	"slices"

	"github.com/grafana/grafana/pkg/services/ngalert/api/tooling/definitions"
	"github.com/grafana/grafana/pkg/services/ngalert/models"
	"github.com/grafana/grafana/pkg/util"
)

func (rev *ConfigRevision) DeleteReceiver(uid string) {
	// Remove the receiver from the configuration.
	rev.Config.AlertmanagerConfig.Receivers = slices.DeleteFunc(rev.Config.AlertmanagerConfig.Receivers, func(r *definitions.PostableApiReceiver) bool {
		return NameToUid(r.GetName()) == uid
	})
}

func (rev *ConfigRevision) CreateReceiver(receiver *models.Receiver) error {
	// Check if the receiver already exists.
	_, err := rev.GetReceiver(receiver.GetUID())
	if err == nil {
		return ErrReceiverExists.Errorf("")
	}
	if !errors.Is(err, ErrReceiverNotFound) {
		return err
	}

	if err := validateAndSetIntegrationUIDs(receiver); err != nil {
		return err
	}

	postable, err := ReceiverToPostableApiReceiver(receiver)
	if err != nil {
		return err
	}

	rev.Config.AlertmanagerConfig.Receivers = append(rev.Config.AlertmanagerConfig.Receivers, postable)

	if err := rev.ValidateReceivers(); err != nil {
		return err
	}

	return nil
}

func (rev *ConfigRevision) UpdateReceiver(receiver *models.Receiver) error {
	existing, err := rev.GetReceiver(receiver.GetUID())
	if err != nil {
		return err
	}

	if err := validateAndSetIntegrationUIDs(receiver); err != nil {
		return err
	}

	postable, err := ReceiverToPostableApiReceiver(receiver)
	if err != nil {
		return err
	}

	// Update receiver in the configuration.
	*existing = *postable

	if err := rev.ValidateReceivers(); err != nil {
		return err
	}

	return nil
}

func (rev *ConfigRevision) ReceiverNameUsedByRoutes(name string) bool {
	return isReceiverInUse(name, []*definitions.Route{rev.Config.AlertmanagerConfig.Route})
}

func (rev *ConfigRevision) GetReceiver(uid string) (*definitions.PostableApiReceiver, error) {
	for _, r := range rev.Config.AlertmanagerConfig.Receivers {
		if NameToUid(r.GetName()) == uid {
			return r, nil
		}
	}
	return nil, ErrReceiverNotFound.Errorf("")
}

func (rev *ConfigRevision) GetReceivers(uids []string) []*definitions.PostableApiReceiver {
	receivers := make([]*definitions.PostableApiReceiver, 0, len(uids))
	for _, r := range rev.Config.AlertmanagerConfig.Receivers {
		if len(uids) == 0 || slices.Contains(uids, NameToUid(r.GetName())) {
			receivers = append(receivers, r)
		}
	}
	return receivers
}

// RenameReceiverInRoutes renames all references to a receiver in routes.
func (rev *ConfigRevision) RenameReceiverInRoutes(oldName, newName string) {
	RenameReceiverInRoute(oldName, newName, rev.Config.AlertmanagerConfig.Route)
}

// ValidateReceivers checks if any receivers have the same Name or if any integrations have the same UID.
func (rev *ConfigRevision) ValidateReceivers() error {
	receivers := make(map[string]struct{}, len(rev.Config.AlertmanagerConfig.Receivers))
	uids := make(map[string]struct{}, len(rev.Config.AlertmanagerConfig.Receivers))
	for _, r := range rev.Config.AlertmanagerConfig.Receivers {
		if _, exists := receivers[r.GetName()]; exists {
			return MakeErrReceiverInvalid(fmt.Errorf("name %q already exists", r.GetName()))
		}
		receivers[r.GetName()] = struct{}{}

		for _, gr := range r.GrafanaManagedReceivers {
			if _, exists := uids[gr.UID]; exists {
				return MakeErrReceiverInvalid(fmt.Errorf("integration with UID %q already exists", gr.UID))
			}
			uids[gr.UID] = struct{}{}
		}
	}
	return nil
}

func RenameReceiverInRoute(oldName, newName string, routes ...*definitions.Route) {
	if len(routes) == 0 {
		return
	}
	for _, route := range routes {
		if route.Receiver == oldName {
			route.Receiver = newName
		}
		RenameReceiverInRoute(oldName, newName, route.Routes...)
	}
}

// isReceiverInUse checks if a receiver is used in a route or any of its sub-routes.
func isReceiverInUse(name string, routes []*definitions.Route) bool {
	if len(routes) == 0 {
		return false
	}
	for _, route := range routes {
		if route.Receiver == name {
			return true
		}
		if isReceiverInUse(name, route.Routes) {
			return true
		}
	}
	return false
}

// validateAndSetIntegrationUIDs validates existing integration UIDs and generates them if they are empty.
func validateAndSetIntegrationUIDs(receiver *models.Receiver) error {
	for _, integration := range receiver.Integrations {
		if integration.UID == "" {
			integration.UID = util.GenerateShortUID()
		} else if err := util.ValidateUID(integration.UID); err != nil {
			return MakeErrReceiverInvalid(fmt.Errorf("integration UID %q is invalid: %w", integration.UID, err))
		}
	}
	return nil
}
