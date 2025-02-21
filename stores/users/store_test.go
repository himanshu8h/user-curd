package users

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
	"user-crud/models"
)

var usr = models.User{
	Id:    1,
	Name:  "himanshu",
	Email: "himanshu8083@gmail.com",
	Phone: "8083860404",
	Age:   25,
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}
func TestUserStore_Insert(t *testing.T) {
	db, mock := NewMock()
	query := "insert into user1(name,email,phone,age) values(?,?,?,?)"

	testcases := []struct {
		desc   string
		inp    models.User
		expOut *models.User
		expErr error
		mock   []interface{}
	}{
		{
			"testcase-1",
			usr,
			&usr,
			nil,
			[]interface{}{
				mock.ExpectExec(query).
					WithArgs(usr.Name, usr.Email, usr.Phone, usr.Age).
					WillReturnResult(sqlmock.NewResult(1, 0)),
			},
		},
		{
			"testcase-2",
			usr,
			&models.User{},
			errors.New("error in executing insert"),
			[]interface{}{
				mock.ExpectPrepare(query).ExpectExec().
					WithArgs(usr.Id, usr.Name, usr.Email, usr.Phone, usr.Age).
					WillReturnError(errors.New("error in executing insert")),
			},
		},
		{
			"testcase-3",
			usr,
			&models.User{},
			errors.New("error in executing insert"),
			[]interface{}{
				mock.ExpectPrepare(query).
					WillReturnError(errors.New("error in executing insert")),
			},
		},
	}
	dbhandler := New(db)
	for _, tcs := range testcases {
		out, err := dbhandler.Insert(&tcs.inp)
		if !reflect.DeepEqual(out, tcs.expOut) {
			t.Errorf("expected %v, got %v", tcs.expOut, out)
		}
		if !reflect.DeepEqual(tcs.expErr, err) {
			t.Errorf("expected error %v, got %v", tcs.expErr, err)
		}
	}
}
func TestUserStore_Update(t *testing.T) {
	db, mock := NewMock()
	dbhandler := New(db)
	testcases := []struct {
		desc   string
		inp    models.User
		expErr error
		mock   []interface{}
	}{
		{
			"testcase-1",
			usr,
			nil,
			[]interface{}{
				mock.ExpectExec("update user1 set name = ?, email = ?, phone = ?, age = ? where id = ?").WithArgs("himanshu", "himanshu8083@gmail.com", "8083860404", 25, 1).
					WillReturnResult(sqlmock.NewResult(0, 1)),
			},
		},
		{
			"testcase-2",
			usr,
			errors.New("error in updating user name"),
			[]interface{}{
				mock.ExpectExec("update user1 set name = ?, email = ?, phone = ?, age = ? where id = ?").WithArgs("himanshu", "himanshu8083@gmail.com", "8083860404", 25, 1).
					WillReturnError(errors.New("error in updating user name")),
			},
		},
	}
	for _, tcs := range testcases {
		_, err := dbhandler.Update(&tcs.inp)
		if !reflect.DeepEqual(err, tcs.expErr) {
			t.Errorf("expect error %v, got %v", tcs.expErr, err)
		}
	}
}
func TestUserStore_Delete(t *testing.T) {
	db, mock := NewMock()
	query := "delete from user1 where id = (?)"
	dbhandler := New(db)
	testcases := []struct {
		desc   string
		inp    int
		expErr error
		mock   []interface{}
	}{
		{
			"testcase-1",
			1,
			nil,
			[]interface{}{
				mock.ExpectExec(query).WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1)),
			},
		},
	}
	for _, tcs := range testcases {
		err := dbhandler.Delete(tcs.inp)

		if !reflect.DeepEqual(err, tcs.expErr) {
			t.Errorf("expected error %v, got %v", tcs.expErr, err)
		}
	}
}
func TestUserStore_GetById(t *testing.T) {
	db, mock := NewMock()
	query := "select * from user1 where id = ?"
	dbhandler := New(db)

	testcases := []struct {
		desc   string
		inp    int
		exp    *models.User
		expErr error
		mock   []interface{}
	}{
		{
			"testcase-1",
			1,
			&usr,
			nil,
			[]interface{}{
				mock.ExpectQuery(query).WithArgs(1).
					WillReturnRows(
						sqlmock.NewRows([]string{"Id", "Name", "Email", "Phone", "Age"}).
							AddRow(
								1,
								"himanshu",
								"himanshu8083@gmail.com",
								"8083860404",
								25,
							),
					),
			},
		},
	}
	for _, tcs := range testcases {
		out, err := dbhandler.GetById(tcs.inp)
		if !reflect.DeepEqual(out, tcs.exp) {
			t.Errorf("expected %v, got %v", tcs.exp, out)
		}
		if !reflect.DeepEqual(tcs.expErr, err) {
			t.Errorf("expected error %v , got %v", tcs.expErr, err)
		}
	}
}
func TestUserStore_GetAll(t *testing.T) {
	db, mock := NewMock()
	query := "select * from user1"
	dbhandler := New(db)

	testcases := []struct {
		desc   string
		inp    int
		exp    []*models.User
		expErr error
		mock   []interface{}
	}{
		{
			"testcase-1",
			1,
			[]*models.User{&usr},
			nil,
			[]interface{}{
				mock.ExpectQuery(query).
					WillReturnRows(
						sqlmock.NewRows([]string{"Id", "Name", "Email", "Phone", "Age"}).
							AddRow(
								1,
								"himanshu",
								"himanshu8083@gmail.com",
								"8083860404",
								25,
							),
					),
			},
		},
		{
			"testcase-2",
			1,
			[]*models.User{},
			errors.New("error in preparing statement"),
			[]interface{}{
				mock.ExpectQuery(query).WillReturnError(errors.New("error in preparing statement")),
			},
		},
	}
	for _, tcs := range testcases {
		out, err := dbhandler.GetAll()
		if !reflect.DeepEqual(out, tcs.exp) {
			t.Errorf("expected %v, got %v", tcs.exp, out)
		}
		if !reflect.DeepEqual(tcs.expErr, err) {
			t.Errorf("expected error %v , got %v", tcs.expErr, err)
		}
	}
}
