package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"forum.christen.net/internal/models"
	"forum.christen.net/internal/validator"
)

type userForm struct {
	Username string `form:"name"`
	Email string `form:"email"`
	Password string	`form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

type threadForm struct {
	Title string `form:"title"`
	validator.Validator
}

type messageForm struct {
	ThreadTitle string `form:"title"`
	Content string `from:"content"`
	validator.Validator
}

// Homepage handler function
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	threads, err := app.threads.ShowLatestThreads()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Threads = threads

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

// Create Account handler function
func (app *application) accountGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userForm{}
	app.render(w, r, http.StatusOK, "account-create.tmpl", data)
}

// Create Account POST handler function
func (app *application) accountPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userForm {
		Username: r.PostForm.Get("username"),
		Email: r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Username), 
			"username", "User name cannot be blank")
	form.CheckField(validator.MaxChars(form.Username, 100), 
			"username", "User name cannot be more than 100 characters long")

	form.CheckField(app.users.ISUniqueEmail(form.Email), 
			"email", "This email is already userd")
	form.CheckField(validator.NotBlank(form.Email), 
			"email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), 
			"email", "This field must be a valid email address")
	form.CheckField(validator.MaxChars(form.Email, 100), 
			"email", "Email cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.Password), 
			"password", "Password cannot be blank")
	form.CheckField(validator.MaxChars(form.Password, 15), 
			"password", "Password cannot be more than 15 characters long")
	form.CheckField(validator.MinChars(form.Password, 8), 
			"password", "Password cannot be less than 8 characters long")
	form.CheckField(validator.HasUppercase(form.Password), 
			"password", "Password must has at least one uppercase letter")
	form.CheckField(validator.HasDigit(form.Password), 
			"password", "Password must has at least one digit")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "account-create.tmpl", data)
		return
	}

	err = app.users.CreateUserAccount(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedEmail) {
            form.AddFieldError("email", "Email address is already in use")
            data := app.newTemplateData(r)
            data.Form = form
            app.render(w, r, http.StatusUnprocessableEntity, "account-create.tmpl", data)
        } else {
            app.serverError(w, r, err)
        }
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Account successfully created!")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// User Login Get handler function
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

// User Login Post handler function
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := userLoginForm{
		Email: r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
    form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email must be a valid email address")
    form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
    if err != nil {
        app.serverError(w, r, err)
        return
    }

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/thread/create", http.StatusSeeOther)
}

// User Logout Post handler function
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Create thread handler function
func (app *application) createThreadGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = threadForm{}
	app.render(w, r, http.StatusOK, "thread-create.tmpl", data)
}

// Create thread POST handler function
func (app *application) createThreadPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userId := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	form := threadForm{
		Title: r.PostForm.Get("title"),
	}

	form.CheckField(validator.NotBlank(form.Title), 
			"title", "This field cannot be blank")
    form.CheckField(validator.MaxChars(form.Title, 100), 
			"title", "This field cannot be more than 100 characters long")
	form.CheckField(app.threads.IfThreadExist(form.Title), 
			"title", "This thread title is already exist")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "thread-create.tmpl", data)
		return
	}

	id, err := app.threads.CreateThread(form.Title, userId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Thread successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/thread/view/%d", id), http.StatusSeeOther)
}

// Display Thread handler function
func (app *application) threadView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	threadMsgs, err := app.threads.ShowThreadMsgs(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	threadMsgs = incrementMsgsId(threadMsgs, id)
	data.Thread = *threadMsgs

	app.render(w, r, http.StatusOK, "thread-view.tmpl", data)
}

func incrementMsgsId(thread *models.Thread, id int) *models.Thread {
	for i := range thread.Messages {
		thread.Messages[i].ThreadId = id
		thread.Messages[i].Id = i + 1
	}
	return thread
}

// A new message handler function
func (app *application) msgGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = messageForm{}

	app.render(w, r, http.StatusOK, "message-create.tmpl", data)
}

// A new message POST handler function
func (app *application) msgPost(w http.ResponseWriter, r *http.Request) {
	var threadId int
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	userId := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	form := messageForm{
		ThreadTitle: r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
	}

	form.CheckField(validator.NotBlank(form.ThreadTitle), 
			"title", "Title cannot be empty")
	form.CheckField(validator.MaxChars(form.ThreadTitle, 100), 
			"title", "Title cannot be more than 100 characters long")

	form.CheckField(validator.NotBlank(form.Content), 
			"content", "The message content cannot be empty")
	form.CheckField(validator.MaxChars(form.Content, 1000), 
			"content", "Message content cannot be more than 1000 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "message-create.tmpl", data)
		return
	}

	if app.threads.IfThreadExist(form.ThreadTitle) {
		threadId, err = app.threads.GetThreadId(form.ThreadTitle)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	} else {
		threadId, err = app.threads.CreateThread(form.ThreadTitle, userId)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	msgId, err := app.msgs.CreateMsg(form.Content, userId, threadId)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Message successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/thread/%d/message/view/%d", threadId, msgId), http.StatusSeeOther)
}

// Account view handler function
func (app *application) msgView(w http.ResponseWriter, r *http.Request) {
	threadId, err := strconv.Atoi(r.PathValue("threadId"))
	if err != nil || threadId < 1 {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg, err := app.msgs.ShowMsg(threadId, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Msg = *msg

	app.render(w, r, http.StatusOK, "message-view.tmpl", data)
}
