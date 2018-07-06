# Simple Web Application


Aplikacja jest solidną bazą i szablonem do tworzenia bardziej skomplikowanych aplikacji webowych. Jej zadaniem jest demonstracja użycia języka Go oraz renderowanych po stronie serwera szablonów html, używających CSS, JavaScript wraz z JQuery.


## Wymagania i uruchamianie

Aplikacja wykorzystuję jedną z trzech dostępnych baz danych:
1. [Bolt](https://github.com/boltdb/bolt)
2. [MySQL](https://www.mongodb.com/)
3. [MongoDB](https://www.mysql.com)

Domyślnym wyborem jest Bolt. Wykorzystanie pozostałych opcji wymaga następującej konfiguracji:
1. **MySql**
    - Wymagane jest posiadanie instancji, przykładowo postawionej na [aws](https://aws.amazon.com/getting-started/tutorials/create-mysql-db/).
    - Następnie w pliku config.json należy uzupełnić poniższe informacje, tak aby pasowały do posiadanej instancji.
        ```
        "MySQL": {
                "Username": "root",
                "Password": "",
                "Name": "gowebapp",
                "Hostname": "127.0.0.1",
                "Port": 3306,
                "Parameter": "?parseTime=true"
        }
        ```

    - W tym samym pliku należy również zmienić:
        ```
        "Type": "Bolt",
        ```
        na
         ```
         "Type": "MySQL",
         ```
2. **MongDB**
    - należy wystartować MongoDB
    - należy ustawić typ bazy w config.json na:
            ```
            "Type": "Bolt",
            ```
            na
             ```
             "Type": "MongoDB",
             ```
      oraz ustawić odpowiednie parametry tak jak w przypadku MySQL.

      ```
      "MongoDB": {
            "URL": "127.0.0.1",
            "Database": "gowebapp"
      },
      ```

Po konfiguracji bazy, włączanie następuje po wywołaniu komendy
    ``` go run gowebapp.go ```, które może wymagać uprawnień administratora.
    Następnie w celu przejścia do panelu logowania należy otworzyć przeglądarkę na: http://localhost


## Struktura

Aplikacja pisana jest zgodnie z wzorcem MVC

Struktura projektu prezentuje się następująco:
```
config                      -- konfiguracja
docs                        -- niniejsza dokumentacja

service                     -- backend aplikacji
service/controller          -- warstwa kontrolera, tj. logika stron i metody GET oraz POST
service/model               -- obiekty bazodanowe oraz zapytania
service/route               -- routing
service/share               -- pomocnicze metody, używane w różnych pakietach

static                      -- lokacja plików js oraz css, jpg oraz czcionek, która jest serwowana przez serwer
template                    -- lokacja szablonów html
```

## Frontend

Frontend zbudowany jest przy renderowanych po stronie serera szablonów html. Technologie, jakie zostały wykorzystane, to bootstrap 3.3.5 oraz jquery. W większości miejsc wykorzystane jest domyślne stylowanie boostrapa.

Ze względu na renderowanie każdej, ze stron po stronie serwera, całość indeksowana jest przez wyszukiwarke google i może zostać łatwo pod tym kątem zooptymalizowana (co byłoby problematyczne, przy pisaniu SPA w którymś z popularnych frameworków).


### HTML Templates

w pakiecie template, znajdują sie szablony html. Są one renderowane po stronie serwera przez silnik szablonów dostępny w standardowej bibliotece języka Go.

Przykładowy pusty szablon, wygląda następująco

```
{{define "title"}}Blank Template{{end}}
{{define "head"}}{{end}}
{{define "content"}}
This is a blank template.
{{end}}
{{define "foot"}}{{end}}
```

W podwójnych klamrach "{" znajdują się polecenia, które zostaną wykonane podczas renderingu. Przykłodowo wstawianie szablonu.

## Backend

Warstwa service zawiera 4 pakiet: controller, model, route oraz share i jest rdzeniem całej aplikacji.

### controller
W pakiecie kontroler znajdują się metody wywoływane bezpośrednio przez router.

W przypadku stron i podstron odpowiada ona za ich renderowanie.

W przypadku zapytań do części backendowej, definiuje ona metody wywoływane bezpośrednio przez router.

### route

Pakiet route implementuje router oraz middleware. Wykorzystano tutaj zewnętrzny pakiet ```github.com/julienschmidt/httprouter``` ze względu na jego dużą prostotę oraz bardzo dobrą optymalizację.

Zaimplementowano 4 middlewary odpowiadająće bezpośrednio za:
1. Autoryzację
2. Ustawianie kontekstu
3. Wypisywanie logów
4. Profiler

Implementowane jest to poprzez zastosowanie dekoratorów zwracających ```HandlerFunc``` lub ```Handler``` tutaj przykładowo

```
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Format("2006-01-02 03:04:05 PM"), r.RemoteAddr, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}
```

Funkcja przyjmuje ```http.Handler```, dekoruje go logowaniem i zwraca taką samą funkcję (Golang pozwala na przypisywanie funkcji do zmiennej).

### model

W pakiecie model zaimplementowane są następujące obiekty bazo-danowe:

```
type Note struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID        uint32        `db:"id" bson:"id,omitempty"`
	Content   string        `db:"content" bson:"content"`
	UserID    bson.ObjectId `bson:"user_id"`
	UID       uint32        `db:"user_id" bson:"userid,omitempty"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
}
```

```
type User struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	ID        uint32        `db:"id" bson:"id,omitempty"`
	FirstName string        `db:"first_name" bson:"first_name"`
	LastName  string        `db:"last_name" bson:"last_name"`
	Email     string        `db:"email" bson:"email"`
	Password  string        `db:"password" bson:"password"`
	StatusID  uint8         `db:"status_id" bson:"status_id"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
}
```

```
type UserStatus struct {
	ID        uint8     `db:"id" bson:"id"`
	Status    string    `db:"status" bson:"status"`
	CreatedAt time.Time `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time `db:"updated_at" bson:"updated_at"`
	Deleted   uint8     `db:"deleted" bson:"deleted"`
}
```

Są one kompatybilne z każdą z baz danych (wyrażenie objęte w "``" odpowiadają za nazwę pola w bazie.

Wspierają podstawowe operacje CRUD

dodatkowo model definiuje następujące błędy:

```
ErrCode = errors.New("Case statement in code is not correct.")
ErrNoResult = errors.New("Result not found.")
ErrUnavailable = errors.New("Database is unavailable.")
ErrUnauthorized = errors.New("User does not have permission to perform this operation.")
```

### shared

#### email

#### passhash

#### recaptcha

#### session

#### view
