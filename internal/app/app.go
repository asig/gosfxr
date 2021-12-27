/*
 * Copyright (c) 2021 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of gosfxr.
 *
 * gosfxr is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * gosfxr is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with gosfxr.  If not, see <http://www.gnu.org/licenses/>.
 */
package app

import (
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"github.com/asig/gosfxr/internal/generator"
	"github.com/asig/gosfxr/internal/ui"
)

type App struct {
	app *gtk.Application
}

func New() *App {

	app := &App{}

	app.app, _ = gtk.ApplicationNew("com.asigner.gosfxr", glib.APPLICATION_FLAGS_NONE)
	app.app.Connect("activate", app.onActivate)
	return app
}

func (a *App) AddWindow(appWindow *ui.AppWindow) {
	a.app.AddWindow(appWindow.GtkWindow())
}

func (a *App) onActivate() {

	g := generator.NewConfig()
	appWindow := ui.NewAppWindow(a.app, g)
	a.AddWindow(appWindow)

	appWindow.Show()
}

func (a *App) Run() {
	a.app.Run(os.Args)

}
