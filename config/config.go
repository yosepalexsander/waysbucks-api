package config

import "os"

var PORT = os.Getenv("PORT")
var CLOUDINARY_URL = os.Getenv("CLOUDINARY_URL")
var DATABASE_URL = os.Getenv("DATABASE_URL")
var GOOGLE_CLIENT_ID = os.Getenv("GOOGLE_CLIENT_ID")
var JWT_SECRET = os.Getenv("JWT_SECRET")
var MIDTRANS_SERVER_KEY = os.Getenv("MIDTRANS_SERVER_KEY")
var MIDTRANS_CLIENT_KEY = os.Getenv("MIDTRANS_CLIENT_KEY")
