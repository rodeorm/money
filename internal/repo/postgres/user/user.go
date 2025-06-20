package user

import (
	"context"
	"fmt"
	"log"
	"crafa/internal/core"
	"crafa/internal/crypt"
	"crafa/internal/logger"
	"time"

	"go.uber.org/zap"
)

// RegUser создает пользователя в БД
func (s *Storage) RegUser(ctx context.Context, u *core.User, domain string) (*core.Session, error) {
	if u.Login == "" { //TODO: проверку на сложность логина
		return nil, fmt.Errorf("логин не может быть пустым")
	}
	if u.Password == "" { //TODO: проверку на сложность пароля
		return nil, fmt.Errorf("пароль не может быть пустым")
	}

	passwordHash, err := crypt.HashPassword(u.Password)
	if err != nil {
		log.Println("RegUser 1", err)
		return nil, err
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("RegUser 2", err)
		return nil, err
	}

	err = tx.Stmtx(s.stmt["insertUser"]).GetContext(ctx, &u.ID, core.RoleReg, u.Login, passwordHash, u.Name, u.FamilyName, u.PatronName, u.Email, u.Phone)

	if err != nil {
		switch err.Error() {
		case `ERROR: duplicate key value violates unique constraint "users_login_key" (SQLSTATE 23505)`:
			tx.Rollback()
			return nil, fmt.Errorf("такой логин уже используется")
		case `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`:
			tx.Rollback()
			return nil, fmt.Errorf("такой адрес электронной почты уже используется")
		}
		log.Println("RegUser 3", err)
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Stmtx(s.stmt["insertMsg"]).ExecContext(ctx, core.MessageTypeConfirm, core.MessageCategoryEmail, u.ID, crypt.GetOneTimePassword(), u.Email)
	if err != nil {
		log.Println("RegUser 4", err)
		tx.Rollback()
		return nil, err
	}

	u.Role = core.Role{ID: core.RoleReg}

	session := &core.Session{User: *u}
	err = tx.Stmtx(s.stmt["insertSession"]).GetContext(ctx, &session.ID, u.ID, time.Now(), time.Now())
	if err != nil {
		log.Println("RegUser 5", err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return session, nil
}

func (s *Storage) ConfirmUserEmail(ctx context.Context, userID int, otp string) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("ConfirmUserEmail 1", err)
		return err
	}
	msg := core.Message{}
	// Проверяем переданный код  UserID = $1 AND Text = $2 AND Email = $3
	err = tx.Stmtx(s.stmt["selectConfMsg"]).GetContext(ctx, &msg, userID, otp)
	if err != nil {
		log.Println("ConfirmUserEmail 2", err)
		tx.Rollback()
		return err
	}

	if msg.ID == 0 {
		tx.Rollback()
		return fmt.Errorf("переданные данные невалидны. нет сообщения для такого пользователя с таким адресом и кодом подтверждения")
	}
	msg.Used = true

	// Обновляем сообщение, что оно было использовано
	_, err = tx.Stmtx(s.stmt["updateMsg"]).ExecContext(ctx, msg.ID, msg.Used, msg.Queued, msg.SendTime)
	if err != nil {
		log.Println("ConfirmUserEmail 2", err)
		tx.Rollback()
		return err
	}

	// Обновляем роль пользователя
	_, err = tx.Stmtx(s.stmt["updateUserRole"]).ExecContext(ctx, userID, core.RoleAuth)
	if err != nil {
		log.Println("ConfirmUserEmail 3", err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Аутентифицирует пользователя на основании данных в БД и возвращает все его данные
func (s *Storage) BaseAuthUser(ctx context.Context, u *core.User) error {

	// Сохраняем введенный пароль
	pass := u.Password

	// Получаем данные пользователя, включая хэш пароля
	err := s.stmt["baseAuthUser"].Get(u, u.Login)
	if err != nil {
		return err
	}
	//Сравнить полученный хэш с паролем
	if !crypt.CheckPasswordHash(pass, u.Password) {
		return fmt.Errorf("некорректный пароль")
	}

	//Добавляем сообщение с одноразовым паролем
	s.stmt["insertMsg"].ExecContext(ctx, core.MessageTypeAuth, core.MessageCategoryEmail, u.ID, crypt.GetOneTimePassword(), u.Email)
	return nil
}

// /AdvAuthUser авторизует пользователя, прошедшего базовую аутентификацию по одноразовому паролю
func (s *Storage) AdvAuthUser(ctx context.Context, u *core.User, otp string, otpLiveTime time.Duration) (*core.Session, error) {
	msg := core.Message{}
	// Проверяем переданный код  UserID = $1 AND Text = $2
	err := s.stmt["selectAuthMsg"].GetContext(ctx, &msg, u.ID, otp)
	if err != nil {
		return nil, fmt.Errorf("неправильный одноразовый пароль")
	}

	if time.Since(msg.SendTime.Time) > otpLiveTime {
		return nil, fmt.Errorf("срок действия одноразового пароля истек")
	}

	// Если всё хорошо с паролем, то получаем данные пользователя, включая роль
	err = s.stmt["selectUser"].GetContext(ctx, u, u.ID)
	if err != nil {
		log.Println("AdvAuthUser 3", err)
		return nil, err
	}

	sn := core.Session{User: *u}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Println("ConfirmUserEmail 1", err)
		return nil, err
	}

	// Затем создаем для него новую сессию и возвращаем её
	err = tx.Stmtx(s.stmt["insertSession"]).GetContext(ctx, &sn.ID, u.ID, time.Now(), time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Обновляем сообщение, проставляем признак, что было использовано
	msg.Used = true

	// Обновляем сообщение, что оно было использовано
	_, err = tx.Stmtx(s.stmt["updateMsg"]).ExecContext(ctx, msg.ID, msg.Used, msg.Queued, msg.SendTime)
	if err != nil {
		log.Println("ConfirmUserEmail 2", err)
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &sn, nil
}

// Возвращает данные пользователя
func (s *Storage) SelectUser(ctx context.Context, u *core.User) error {
	err := s.stmt["selectUser"].GetContext(ctx, u, u.ID)
	return err
}

// Возвращает данные всех пользователей
func (s *Storage) SelectAllUsers(ctx context.Context) ([]core.User, error) {
	u := make([]core.User, 0)
	err := s.stmt["selectAllUsers"].SelectContext(ctx, &u)
	if err != nil {
		logger.Log.Info("SelectAllUsers",
			zap.Error(err))
		return nil, err
	}
	return u, nil
}

// Обновляет данные пользователя
func (s *Storage) UpdateUser(ctx context.Context, u *core.User) error {
	//roleid = $2, login = $3, name = $4, familyname = $5, patronname = $6, email = $7, phone = $8
	_, err := s.stmt["updateUser"].ExecContext(ctx, u.ID, u.Role.ID, u.Login, u.Name, u.FamilyName, u.PatronName, u.Email, u.Phone)
	if err != nil {
		return err
	}

	return nil
}

// Меняет пароль пользователю
func (s *Storage) ChangeUserPassword(ctx context.Context, u *core.User) error {
	// Хэшируем пароль
	pwdHash, err := crypt.HashPassword(u.Password)
	if err != nil {
		return err
	}
	_, err = s.stmt["changeUserPassword"].ExecContext(ctx, u.ID, pwdHash)
	if err != nil {
		return err
	}
	return nil
}
