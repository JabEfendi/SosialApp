package controllers

import (
	"backend/db"
	"backend/models"
	"net/http"
    "bytes"
	"encoding/json"
	"io/ioutil"
	"time"
    "backend/firebase"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)