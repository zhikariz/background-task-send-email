package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func main() {
	cfg, err := NewConfig()
	checkError(err)

	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port))
		},
	}

	pool := work.NewWorkerPool(Context{}, 10, cfg.Namespace, redisPool)

	pool.Middleware((*Context).Log)
	pool.Middleware((*Context).FindCustomer)

	pool.Job("send_welcome_email", (*Context).SendWelcomeEmail)

	pool.Start()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	pool.Stop()
}

type Context struct {
	email  string
	userID int64
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting Job: ", job.Name)
	return next()
}

func (c *Context) FindCustomer(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Finding customer: ", c.email)
	if _, ok := job.Args["user_id"]; !ok {
		c.userID = job.ArgInt64("user_id")
		c.email = job.ArgString("email_address")
		if err := job.ArgError(); err != nil {
			return fmt.Errorf("arg error %v", err.Error())
		}
	}
	return next()
}

func (c *Context) SendWelcomeEmail(job *work.Job) error {
	cfg, err := NewConfig()
	checkError(err)
	addr := job.ArgString("email_address")
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Sending welcome email to: ", addr)

	from, password := "info@amygdala.cloud", cfg.SMTP.Password

	to := []string{
		addr,
	}

	smtpHost := cfg.SMTP.Host
	// smtpPort := cfg.SMTP.Port
	message := "This is a welcome email message."

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", "Welcome Email")
	mailer.SetBody("text/plain", message)

	dialer := gomail.NewDialer(smtpHost, 587, from, password)

	err = dialer.DialAndSend(mailer)
	checkError(err)

	fmt.Println("Email sent successfully")
	return nil
}

type Config struct {
	Namespace string
	Redis     RedisConfig
	SMTP      SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     string
	Password string
}

type RedisConfig struct {
	Host string
	Port string
}

func NewConfig() (*Config, error) {
	cfg := new(Config)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	checkError(err)
	err = viper.Unmarshal(&cfg)
	checkError(err)

	return cfg, nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Context) SendTicketPaid(job *work.Job) error {
	cfg, err := NewConfig()
	checkError(err)
	addr := job.ArgString("email_address")
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Sending welcome email to: ", addr)

	from, password := "info@amygdala.cloud", cfg.SMTP.Password

	to := []string{
		addr,
	}

	smtpHost := cfg.SMTP.Host
	// smtpPort := cfg.SMTP.Port
	message := "This is a welcome email message."

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", "Welcome Email")
	mailer.SetBody("text/plain", message)

	dialer := gomail.NewDialer(smtpHost, 587, from, password)

	err = dialer.DialAndSend(mailer)
	checkError(err)

	fmt.Println("Email sent successfully")
	return nil
}
