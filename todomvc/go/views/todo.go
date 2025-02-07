package views

import (
	"github.com/go-humble/examples/todomvc/go/models"
	"github.com/go-humble/examples/todomvc/go/templates"
	"github.com/go-humble/temple/temple"
	"github.com/go-humble/view"
	"honnef.co/go/js/dom"
)

var (
	todoTmpl = templates.MustGetPartial("todo")
)

// Todo is a view for a single todo item.
type Todo struct {
	Model *models.Todo
	tmpl  *temple.Partial
	view.DefaultView
}

// NewTodo creates and returns a new Todo view, using the given todo as the
// model.
func NewTodo(todo *models.Todo) *Todo {
	return &Todo{
		Model: todo,
		tmpl:  todoTmpl,
	}
}

// Render renders the Todo view and satisfies the view.View interface.
func (v *Todo) Render() error {
	if err := v.tmpl.ExecuteEl(v.Element(), v.Model); err != nil {
		return err
	}
	v.delegateEvents()
	return nil
}

// delegateEvents adds all the needed event listeners to the Todo view.
func (v *Todo) delegateEvents() {
	view.AddEventListener(v, "click", ".toggle", v.Toggle)
	view.AddEventListener(v, "click", ".destroy", v.Remove)
	view.AddEventListener(v, "dblclick", "label", v.Edit)
	view.AddEventListener(v, "blur", ".edit", v.CommitEdit)
	view.AddEventListener(v, "keypress", ".edit",
		triggerOnKeyCode(enterKey, v.CommitEdit))
	view.AddEventListener(v, "keydown", ".edit",
		triggerOnKeyCode(escapeKey, v.CancelEdit))

}

// Toggle toggles the completeness of the todo.
func (v *Todo) Toggle(ev dom.Event) {
	v.Model.Toggle()
}

// Remove removes the todo form the list.
func (v *Todo) Remove(ev dom.Event) {
	v.Model.Remove()
}

// Edit puts the Todo view into an editing state, changing it's appearance and
// allowing it to be edited.
func (v *Todo) Edit(ev dom.Event) {
	li := v.Element().QuerySelector("li")
	addClass(li, "editing")
	input, ok := li.QuerySelector(".edit").(*dom.HTMLInputElement)
	if !ok {
		panic("Could not convert to dom.HTMLInputElement")
	}
	input.Focus()
	// Move the cursor to the end of the input.
	input.SelectionStart = input.SelectionEnd + len(input.Value)
}

// CommitEdit sets the title of the todo to the new title. After the edit has
// been committed, the todo is no longer in the editing state.
func (v *Todo) CommitEdit(ev dom.Event) {
	li := v.Element().QuerySelector("li")
	input, ok := li.QuerySelector(".edit").(*dom.HTMLInputElement)
	if !ok {
		panic("Could not convert to dom.HTMLInputElement")
	}
	v.Model.SetTitle(input.Value)
}

// CancelEdit resets the title of the todo to its old value. It does not commit
// the edit. After the edit has been canceled, the todo is no longer in the
// editing state.
func (v *Todo) CancelEdit(ev dom.Event) {
	li := v.Element().QuerySelector("li")
	removeClass(li, "editing")
	input, ok := li.QuerySelector(".edit").(*dom.HTMLInputElement)
	if !ok {
		panic("Could not convert to dom.HTMLInputElement")
	}
	input.Value = v.Model.Title()
	input.Blur()
}
