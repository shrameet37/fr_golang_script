import psycopg2

global connection

def openDataBaseConnection(env):

    try:
        if env == 'dev1':
            connection = psycopg2.connect(user="diagnostics_dev1_role",
                                  password="r1W8YH2E6mlwVLoD5SwM",
                                  host="localhost",
                                  port="5445",
                                  database="diagnostics_dev1")
        elif env == 'prod':
            connection = psycopg2.connect(user="facerecognition_management_v3_prod_role",
                                  password="YTkxY2Y4NGFiZTZj",
                                  host="localhost",
                                  port="5433",
                                  database="facerecognition_management_v3_prod")
            
        elif env == 'usProd':
            connection = psycopg2.connect(user="diagnostics_usprod_role",
                                  password="t286n9D7Z6D3XDRRD6MAsQHJ8tP3e8",
                                  host="localhost",
                                  port="5465",
                                  database="diagnostics_usprod")
        else:
            return None,None
            

        cursor = connection.cursor()

        print("Database connection opened!")

        return connection,cursor

    except Exception as e:
        errorMessage = str(e)
        print("Something went wrong when opening conenction!!!")
        print(errorMessage)

def closeDataBaseConnection(connection,cursor):

    try:

        connection.commit()
        cursor.close()
        connection.close()

        print("Database connection closed!")

    except Exception as e:
        errorMessage = str(e)
        print("Something went wrong when closing connection!!!")
        print(errorMessage)

def closeDataBaseConnectionWithoutCommit(connection,cursor):

    try:

        cursor.close()
        connection.close()

        print("Database connection closed Without Commit!")

    except Exception as e:
        errorMessage = str(e)
        print("Something went wrong when closing connection!!!")
        print(errorMessage)
