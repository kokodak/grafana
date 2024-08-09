package ngalert

import (
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/dashboards"
	"github.com/grafana/grafana/pkg/services/datasources"
	ac "github.com/grafana/grafana/pkg/services/ngalert/accesscontrol"
	"github.com/grafana/grafana/pkg/services/org"
)

const AlertRolesGroup = "Alerting"

var (
	rulesReaderRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.rules:reader",
			DisplayName: "Rules Reader",
			Description: "Read alert rules in all Grafana folders and external providers",
			Group:       AlertRolesGroup,
			Permissions: []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingRuleRead,
					Scope:  dashboards.ScopeFoldersAll,
				},
				{
					Action: accesscontrol.ActionAlertingRuleExternalRead,
					Scope:  datasources.ScopeAll,
				},
				{
					Action: accesscontrol.ActionAlertingSilencesRead,
					Scope:  dashboards.ScopeFoldersAll,
				},
				// Following are needed for simplified notification policies
				{
					Action: accesscontrol.ActionAlertingNotificationsTimeIntervalsRead,
				},
				{
					Action: accesscontrol.ActionAlertingReceiversList,
				},
			},
		},
	}

	rulesWriterRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.rules:writer",
			DisplayName: "Rules Writer",
			Description: "Add, update, and delete rules in any Grafana folder and external providers",
			Group:       AlertRolesGroup,
			Permissions: accesscontrol.ConcatPermissions(rulesReaderRole.Role.Permissions, []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingRuleCreate,
					Scope:  dashboards.ScopeFoldersAll,
				},
				{
					Action: accesscontrol.ActionAlertingRuleUpdate,
					Scope:  dashboards.ScopeFoldersAll,
				},
				{
					Action: accesscontrol.ActionAlertingRuleDelete,
					Scope:  dashboards.ScopeFoldersAll,
				},
				{
					Action: accesscontrol.ActionAlertingRuleExternalWrite,
					Scope:  datasources.ScopeAll,
				},
				{
					Action: accesscontrol.ActionAlertingSilencesWrite,
					Scope:  dashboards.ScopeFoldersAll,
				},
				{
					Action: accesscontrol.ActionAlertingSilencesCreate,
					Scope:  dashboards.ScopeFoldersAll,
				},
			}),
		},
	}

	instancesReaderRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.instances:reader",
			DisplayName: "Instances and Silences Reader",
			Description: "Read instances and silences of Grafana and external providers",
			Group:       AlertRolesGroup,
			Permissions: []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingInstanceRead,
				},
				{
					Action: accesscontrol.ActionAlertingInstancesExternalRead,
					Scope:  datasources.ScopeAll,
				},
			},
		},
	}

	instancesWriterRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.instances:writer",
			DisplayName: "Silences Writer",
			Description: "Add and update silences in Grafana and external providers",
			Group:       AlertRolesGroup,
			Permissions: accesscontrol.ConcatPermissions(instancesReaderRole.Role.Permissions, []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingInstanceCreate,
				},
				{
					Action: accesscontrol.ActionAlertingInstanceUpdate,
				},
				{
					Action: accesscontrol.ActionAlertingInstancesExternalWrite,
					Scope:  datasources.ScopeAll,
				},
			}),
		},
	}

	notificationsReaderRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.notifications:reader",
			DisplayName: "Notifications Reader",
			Description: "Read notification policies and contact points in Grafana and external providers",
			Group:       AlertRolesGroup,
			Permissions: []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingNotificationsRead,
				},
				{
					Action: accesscontrol.ActionAlertingNotificationsExternalRead,
					Scope:  datasources.ScopeAll,
				},
				{
					Action: accesscontrol.ActionAlertingNotificationsTimeIntervalsRead,
				},
				{
					Action: accesscontrol.ActionAlertingReceiversRead,
					Scope:  ac.ScopeReceiversAll,
				},
			},
		},
	}

	notificationsWriterRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.notifications:writer",
			DisplayName: "Notifications Writer",
			Description: "Add, update, and delete contact points and notification policies in Grafana and external providers",
			Group:       AlertRolesGroup,
			Permissions: accesscontrol.ConcatPermissions(notificationsReaderRole.Role.Permissions, []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingNotificationsWrite,
				},
				{
					Action: accesscontrol.ActionAlertingNotificationsExternalWrite,
					Scope:  datasources.ScopeAll,
				},
			}),
		},
	}

	alertingReaderRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting:reader",
			DisplayName: "Full read-only access",
			Description: "Read alert rules, instances, silences, contact points, and notification policies in Grafana and all external providers",
			Group:       AlertRolesGroup,
			Permissions: accesscontrol.ConcatPermissions(rulesReaderRole.Role.Permissions, instancesReaderRole.Role.Permissions, notificationsReaderRole.Role.Permissions),
		},
		Grants: []string{string(org.RoleViewer)},
	}

	alertingWriterRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting:writer",
			DisplayName: "Full access",
			Description: "Add,update and delete alert rules, instances, silences, contact points, and notification policies in Grafana and all external providers",
			Group:       AlertRolesGroup,
			Permissions: accesscontrol.ConcatPermissions(rulesWriterRole.Role.Permissions, instancesWriterRole.Role.Permissions, notificationsWriterRole.Role.Permissions),
		},
		Grants: []string{string(org.RoleEditor), string(org.RoleAdmin)},
	}

	alertingProvisionerRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.provisioning:writer",
			DisplayName: "Access to alert rules provisioning API",
			Description: "Manage all alert rules, contact points, notification policies, silences, etc. in the organization via provisioning API.",
			Group:       AlertRolesGroup,
			Permissions: []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingProvisioningRead, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingProvisioningWrite, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingRulesProvisioningRead, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingRulesProvisioningWrite, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingNotificationsProvisioningRead, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingNotificationsProvisioningWrite, // organization scope
				},
				{
					Action: dashboards.ActionFoldersRead,
					Scope:  dashboards.ScopeFoldersAll,
				},
			},
		},
		Grants: []string{string(org.RoleAdmin)},
	}

	alertingProvisioningReaderWithSecretsRole = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.provisioning.secrets:reader",
			DisplayName: "Read via Provisioning API + Export Secrets",
			Description: "Read all alert rules, contact points, notification policies, silences, etc. in the organization via provisioning API and use export with decrypted secrets",
			Group:       AlertRolesGroup,
			Permissions: []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingProvisioningReadSecrets, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingProvisioningRead, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingRulesProvisioningRead, // organization scope
				},
				{
					Action: accesscontrol.ActionAlertingNotificationsProvisioningRead, // organization scope
				},
			},
		},
		Grants: []string{string(org.RoleAdmin)},
	}

	alertingProvisioningStatus = accesscontrol.RoleRegistration{
		Role: accesscontrol.RoleDTO{
			Name:        accesscontrol.FixedRolePrefix + "alerting.provisioning.provenance:writer",
			DisplayName: "Set provisioning status",
			Description: "Set provisioning status for alerting resources. Should be used together with other regular roles (Notifications Writer and/or Rules Writer)",
			Group:       AlertRolesGroup,
			Permissions: []accesscontrol.Permission{
				{
					Action: accesscontrol.ActionAlertingProvisioningSetStatus, // organization scope
				},
			},
		},
		Grants: []string{string(org.RoleAdmin), string(org.RoleEditor)},
	}
)

func DeclareFixedRoles(service accesscontrol.Service) error {
	return service.DeclareFixedRoles(
		rulesReaderRole, rulesWriterRole,
		instancesReaderRole, instancesWriterRole,
		notificationsReaderRole, notificationsWriterRole,
		alertingReaderRole, alertingWriterRole, alertingProvisionerRole, alertingProvisioningReaderWithSecretsRole, alertingProvisioningStatus,
	)
}
