// Тесты хранилища.
// Внимание! Тест почистит базу данных для работы без ошибок
package storage_test

import (
	"testing"

	// Пакет gofakeit генерирует реалистичные данные, такие как username(в тесте как login) и пароль
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/uncomonq/FinalCalc/internal/storage"
	"github.com/uncomonq/FinalCalc/pckg/consts"
	"github.com/uncomonq/FinalCalc/pckg/types"
)

// Простой тест для проверки открытия/закрытия
func TestStorageSimple(t *testing.T) {
	storage, err := storage.Open(consts.TestStoragePath)
	if err != nil {
		t.Fatalf("Failed to open storage: %v", err)
	}
	err = storage.Close()
	if err != nil {
		t.Fatalf("Failed to close storage: %v", err)
	}
}

// Тест на работу с пользователями
func TestStorageWorkUsers(t *testing.T) {
	storage, err := storage.Open(consts.TestStoragePath) // Открываем базу данных
	if err != nil {
		t.Fatalf("Failed to open storage: %v", err) // Ошибка открытия базы данных
	}
	testcases := []struct {
		Name  string
		Users []*types.User
	}{
		{
			Name: "one",
			Users: []*types.User{
				{
					ID:       uuid.NewString(),                                     // генерируем id
					Login:    gofakeit.Username(),                                  // генерируем имя пользователя
					Password: gofakeit.Password(true, true, true, false, false, 8), // генерируем пароль
				},
			},
		},
		{
			Name: "two",
			Users: []*types.User{
				{
					ID:       uuid.NewString(),
					Login:    gofakeit.Username(),
					Password: gofakeit.Password(true, true, true, false, false, 8),
				},
				{
					ID:       uuid.NewString(),
					Login:    gofakeit.Username(),
					Password: gofakeit.Password(true, true, true, false, false, 8),
				},
			},
		},
		{
			Name: "three",
			Users: []*types.User{
				{
					ID:       uuid.NewString(),
					Login:    gofakeit.Username(),
					Password: gofakeit.Password(true, true, true, false, false, 8),
				},
				{
					ID:       uuid.NewString(),
					Login:    gofakeit.Username(),
					Password: gofakeit.Password(true, true, true, false, false, 8),
				},
				{
					ID:       uuid.NewString(),
					Login:    gofakeit.Username(),
					Password: gofakeit.Password(true, true, true, false, false, 8),
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			err = storage.Clear() // перед началом теста очищаем базу данных
			if err != nil {
				t.Fatalf("Failed to clear storage: %v", err)
			}
			for i, u := range tc.Users { // добавляем пользователей
				err := storage.InsertUser(u)
				if err != nil {
					t.Fatalf("Failed to insert user #%d: %v", i, err)
				}
			}

			selectUsers, err := storage.SelectAllUsers() // получаем список, проверяем длину
			if err != nil {
				t.Fatalf("Failed to select users: %v", err)
			}
			if len(selectUsers) != len(tc.Users) {
				t.Fatalf("invalid data length: expected %d, but got: %d", len(tc.Users), len(selectUsers))
			}

			for i, eu := range tc.Users { // проверяем содержимое списка
				var su *types.User // порядок не сохраняется, ищем по id
				for _, u := range selectUsers {
					if u.ID == eu.ID {
						su = u
						break
					}
				}

				if su == nil {
					t.Fatalf("User with id %s not found", eu.ID)
				}

				if su.ID != eu.ID {
					t.Fatalf("Selected user #%d: invalid ID: expected: %s, but got: %s", i, eu.ID, su.ID)
				}
				if su.Login != eu.Login {
					t.Fatalf("Selected user #%d: invalid login: expected: %s, but got: %s", i, eu.Login, su.Login)
				}
				if su.Password != eu.Password {
					t.Fatalf("Selected user #%d: invalid password: expected: %s, but got: %s", i, eu.Password, su.Password)
				}
			}
		})
	}
	err = storage.Close() // Закрываем базу данных
	if err != nil {
		t.Fatalf("Failed to close storage: %v", err)
	}
}

// Тест на работу с выражениями
func TestStorageWorkExpressions(t *testing.T) {
	storage, err := storage.Open(consts.TestStoragePath) // Открываем базу данных
	if err != nil {
		t.Fatalf("Failed to open storage: %v", err) // Ошибка открытия базы данных
	}
	testcases := []struct {
		Name        string
		Expressions []*types.ExpressionWithID
	}{
		{
			Name: "one",
			Expressions: []*types.ExpressionWithID{
				{
					ID: gofakeit.Uint32(),
					Expression: types.Expression{
						Data:   "2+2", // какие-то данные
						Status: "OK",
						Result: 4,
					},
				},
			},
		},
		{
			Name: "two",
			Expressions: []*types.ExpressionWithID{
				{
					ID: gofakeit.Uint32(),
					Expression: types.Expression{
						Data:   "2+2",
						Status: "OK",
						Result: 4,
					},
				},
				{
					ID: gofakeit.Uint32(),
					Expression: types.Expression{
						Data:   "invalid",
						Status: "error",
						Result: 0,
					},
				},
			},
		},
		{
			Name: "three",
			Expressions: []*types.ExpressionWithID{
				{
					ID: gofakeit.Uint32(),
					Expression: types.Expression{
						Data:   "2+2",
						Status: "OK",
						Result: 4,
					},
				},
				{
					ID: gofakeit.Uint32(),
					Expression: types.Expression{
						Data:   "invalid",
						Status: "error",
					},
				},
				{
					ID: gofakeit.Uint32(),
					Expression: types.Expression{
						Data:   "2+(2/20)*100",
						Status: "Wait",
					},
				},
			},
		},
	}
	testUser := &types.User{
		ID:       uuid.NewString(),
		Login:    gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, false, false, 8),
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			err = storage.Clear() // перед началом теста очищаем базу данных
			if err != nil {
				t.Fatalf("Failed to clear storage: %v", err)
			}

			for i, u := range tc.Expressions { // добавляем все выражения
				err := storage.InsertExpression(u, testUser)
				if err != nil {
					t.Fatalf("Failed to insert expression #%d: %v", i, err)
				}
			}

			// Проверяем как выражения для это пользователя, так и все выражения

			// Получаем списки, проверяем длину
			selectExpressions, err := storage.SelectExpressionsForUser(testUser)
			if err != nil {
				t.Fatalf("Failed to select expressions: %v", err)
			}
			if len(selectExpressions) != len(tc.Expressions) {
				t.Fatalf("invalid expression length: expected %d, but got: %d", len(tc.Expressions), len(selectExpressions))
			}
			selectExpressionsAll, err := storage.SelectExpressions()

			if err != nil {
				t.Fatalf("Failed to select all expressions: %v", err)
			}
			if len(selectExpressionsAll) != len(tc.Expressions) {
				t.Fatalf("invalid all expression length: expected %d, but got: %d", len(tc.Expressions), len(selectExpressionsAll))
			}

			for i, ee := range tc.Expressions { // проверяем содержимое списка
				var se *types.ExpressionWithID // порядок не сохраняется, ищем по id
				for _, e := range selectExpressions {
					if e.ID == ee.ID {
						se = e
						break
					}
				}

				if se == nil {
					t.Fatalf("Expression with id #%d not found", ee.ID)
				}

				if se.ID != ee.ID {
					t.Fatalf("Selected expression #%d: invalid ID: expected: %d, but got: %d", i, ee.ID, se.ID)
				}
				if se.Data != ee.Data {
					t.Fatalf("Selected expression #%d: invalid data: expected: %s, but got: %s", i, ee.Data, se.Data)
				}
				if se.Status != ee.Status {
					t.Fatalf("Selected expression #%d: invalid status: expected: %s, but got: %s", i, ee.Status, se.Status)
				}
				if se.Result != ee.Result {
					t.Fatalf("Selected expression #%d: invalid result: expected: %f, but got: %f", i, ee.Result, se.Result)
				}

				var sea *types.ExpressionWithID // порядок не сохраняется, ищем по id
				for _, e := range selectExpressions {
					if e.ID == ee.ID {
						sea = e
						break
					}
				}

				if sea == nil {
					t.Fatalf("Expression with id #%d not found in all list", ee.ID)
				}

				if sea.ID != ee.ID {
					t.Fatalf("Selected expression in all list #%d: invalid ID: expected: %d, but got: %d", i, ee.ID, sea.ID)
				}
				if sea.Data != ee.Data {
					t.Fatalf("Selected expression in all list #%d: invalid data: expected: %s, but got: %s", i, ee.Data, sea.Data)
				}
				if sea.Status != ee.Status {
					t.Fatalf("Selected expression in all list #%d: invalid status: expected: %s, but got: %s", i, ee.Status, sea.Status)
				}
				if sea.Result != ee.Result {
					t.Fatalf("Selected expression in all list #%d: invalid result: expected: %f, but got: %f", i, ee.Result, sea.Result)
				}
			}
		})
	}
	err = storage.Close() // Закрываем базу данных
	if err != nil {
		t.Fatalf("Failed to close storage: %v", err)
	}
}
