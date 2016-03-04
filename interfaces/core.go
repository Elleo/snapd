// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package interfaces

import (
	"errors"
	"fmt"
	"regexp"
)

// Plug represents the potential of a given snap to connect to a slot.
type Plug struct {
	Snap        string                 `json:"snap"`
	Name        string                 `json:"plug"`
	Interface   string                 `json:"interface"`
	Attrs       map[string]interface{} `json:"attrs,omitempty"`
	Apps        []string               `json:"apps,omitempty"`
	Label       string                 `json:"label"`
	Connections []SlotRef              `json:"connections,omitempty"`
}

// PlugRef is a reference to a plug.
type PlugRef struct {
	Snap string `json:"snap"`
	Name string `json:"plug"`
}

// Slot represents a capacity offered by a snap.
type Slot struct {
	Snap        string                 `json:"snap"`
	Name        string                 `json:"slot"`
	Interface   string                 `json:"interface"`
	Attrs       map[string]interface{} `json:"attrs,omitempty"`
	Apps        []string               `json:"apps,omitempty"`
	Label       string                 `json:"label"`
	Connections []PlugRef              `json:"connections,omitempty"`
}

// SlotRef is a reference to a slot.
type SlotRef struct {
	Snap string `json:"snap"`
	Name string `json:"slot"`
}

// Interfaces holds information about a list of plugs and slots, and their connections.
type Interfaces struct {
	Plugs []*Plug `json:"plugs"`
	Slots []*Slot `json:"slots"`
}

// Interface describes a group of interchangeable capabilities with common features.
// Interfaces act as a contract between system builders, application developers
// and end users.
type Interface interface {
	// Unique and public name of this interface.
	Name() string

	// SanitizePlug checks if a plug is correct, altering if necessary.
	SanitizePlug(plug *Plug) error

	// SanitizeSlot checks if a slot is correct, altering if necessary.
	SanitizeSlot(slot *Slot) error

	// PermanentPlugSnippet returns permanent, plug-side security snippet.
	//
	// Permanent security snippet can be used to grant permissions to a snap that
	// has a plug of a given interface even before the plug is connected to a
	// slot.
	//
	// An empty snippet is returned when there are no additional permissions
	// that are required to implement this interface. ErrUnknownSecurity error
	// is returned when the plug cannot deal with the requested security
	// system.
	PermanentPlugSnippet(plug *Plug, securitySystem SecuritySystem) ([]byte, error)

	// PlugSnippet returns connection-specific, plug-side security snippet.
	//
	// Connection-specific security snippet can be used to grant permission to
	// a snap that has a plug of a given interface connected to a slot in
	// another snap.
	//
	// The snippet should be specific to both the plug and the slot. If the
	// slot is not necessary then consider using PermanentPlugSnippet()
	// instead.
	//
	// An empty snippet is returned when there are no additional permissions
	// that are required to implement this interface. ErrUnknownSecurity error
	// is returned when the plug cannot deal with the requested security
	// system.
	PlugSnippet(plug *Plug, slot *Slot, securitySystem SecuritySystem) ([]byte, error)

	// PermanentSlotSnippet returns permanent, slot-side security snippet.
	//
	// Permanent security snippet can be used to grant permissions to a snap that
	// has a slot of a given interface even before the first connection to that
	// slot is made.
	//
	// An empty snippet is returned when there are no additional permissions
	// that are required to implement this interface. ErrUnknownSecurity error
	// is returned when the plug cannot deal with the requested security
	// system.
	PermanentSlotSnippet(slot *Slot, securitySystem SecuritySystem) ([]byte, error)

	// SlotSnippet returns connection-specific, slot-side security snippet.
	//
	// Connection-specific security snippet can be used to grant permission to
	// a snap that has a slot of a given interface connected to a plug in
	// another snap.
	//
	// The snippet should be specific to both the plug and the slot, if the
	// plug is not necessary then consider using PermanentSlotSnippet()
	// instead.
	//
	// An empty snippet is returned when there are no additional permissions
	// that are required to implement this interface. ErrUnknownSecurity error
	// is returned when the plug cannot deal with the requested security
	// system.
	SlotSnippet(plug *Plug, slot *Slot, securitySystem SecuritySystem) ([]byte, error)
}

// SecuritySystem is a name of a security system.
type SecuritySystem string

const (
	// SecurityAppArmor identifies the apparmor security system.
	SecurityAppArmor SecuritySystem = "apparmor"
	// SecuritySecComp identifies the seccomp security system.
	SecuritySecComp SecuritySystem = "seccomp"
	// SecurityDBus identifies the DBus security system.
	SecurityDBus SecuritySystem = "dbus"
	// SecurityUDev identifies the UDev security system.
	SecurityUDev SecuritySystem = "udev"
)

var (
	// ErrUnknownSecurity is reported when a interface is unable to deal with a given security system.
	ErrUnknownSecurity = errors.New("unknown security system")
)

// Regular expression describing correct identifiers.
var validName = regexp.MustCompile("^[a-z](?:-?[a-z0-9])*$")

// ValidateName checks if a string can be used as a plug or slot name.
func ValidateName(name string) error {
	valid := validName.MatchString(name)
	if !valid {
		return fmt.Errorf("invalid interface name: %q", name)
	}
	return nil
}
