package server

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a standard successful API response.
type Response struct {
	Message string `json:"message,omitempty"`
	Status  string `json:"status,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse represents a standard error API response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// ResponseOk (200 OK) indicates that the request has succeeded.
//
// The payload sent in a 200 response depends on the request method.
// Common use: successful GET, PUT, or POST where the result is returned.
func (c *Ctx) ResponseOk(data any) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

// ResponseCreated (201 Created) indicates that the request has been fulfilled
// and has resulted in one or more new resources being created.
//
// The primary resource created is typically identified by a Location header
// field or by the URI of the request.
func (c *Ctx) ResponseCreated(data any) error {
	if data == nil {
		return c.SendStatus(fiber.StatusCreated)
	}
	return c.Status(fiber.StatusCreated).JSON(data)
}

// ResponseNoContent (204 No Content) indicates that the server has successfully
// fulfilled the request and that there is no additional content to send in the
// response payload body.
//
// Common use: successful DELETE or PUT where no data needs to be returned.
func (c *Ctx) ResponseNoContent() error {
	return c.SendStatus(fiber.StatusNoContent)
}

// ResponseBadRequest (400 Bad Request) indicates that the server cannot or
// will not process the request due to something that is perceived to be a
// client error (e.g., malformed request syntax, invalid request message
// framing, or deceptive request routing).
func (c *Ctx) ResponseBadRequest(message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": message,
	})
}

// ResponseUnauthorized (401 Unauthorized) indicates that the request has not
// been applied because it lacks valid authentication credentials for the
// target resource.
//
// This status code is sent with a WWW-Authenticate header field that contains
// information on how to authorize correctly.
func (c *Ctx) ResponseUnauthorized(message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": message,
	})
}

// ResponseForbidden (403 Forbidden) indicates that the server understood the
// request but refuses to authorize it.
//
// Unlike 401, providing credentials will not make a difference. The access
// is permanently forbidden for the current user/context.
func (c *Ctx) ResponseForbidden(message string) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"error": message,
	})
}

// ResponseNotFound (404 Not Found) indicates that the origin server did not
// find a current representation for the target resource or is not willing to
// disclose that one exists.
func (c *Ctx) ResponseNotFound(message string) error {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": message,
	})
}

// ResponseConflict (409 Conflict) indicates that the request could not be
// processed because of conflict in the current state of the resource, such as
// an edit conflict between multiple simultaneous updates.
func (c *Ctx) ResponseConflict(message string) error {
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		"error": message,
	})
}

// ResponseInternalServerError (500 Internal Server Error) indicates that the
// server encountered an unexpected condition that prevented it from
// fulfilling the request.
//
// This is the "catch-all" error response for server-side exceptions.
func (c *Ctx) ResponseInternalServerError(message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": message,
	})
}
