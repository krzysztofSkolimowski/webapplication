# Simple Web Application

## Wstęp

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

Po konfiguracji bazy, włączanie następuje po wywołaniu komendy
    ``` go run gowebapp.go ```, które może wymagać uprawnień administratora.
    Następnie w celu przejścia do panelu logowania należy otworzyć przeglądarkę na: http://localhost



## Struktura

## Frontend

### HTML Templates

## Backend

### controller

### model

### route

### shared

#### database

#### email

#### jsonconfig

#### passhash

#### racaptcha

#### session

## Technologie