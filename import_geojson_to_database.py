import json
import mysql.connector
from mysql.connector import Error
import getpass
import sys
import geojson


# Function to create database connection
 
def create_connection(host_name, user_name, user_password, db_name):
    connection = None
    try:
        connection = mysql.connector.connect(
            host=host_name,
            user=user_name,
            passwd=user_password,
            database=db_name
        )
        print("Connection to MySQL DB successful")
    except Error as e:
        print(f"The error '{e}' occurred")
    return connection

# Function to execute query
def execute_query(connection, query):
    cursor = connection.cursor()
    try:
        cursor.execute(query)
        connection.commit()
        print("Query executed successfully")
    except Error as e:
        print(f"The error '{e}' occurred")

# Read GeoJSON file
with open('data_harvest/deduped_data.json') as file:
    geojson_data = json.load(file)

#take password as CLI
db_password = getpass.getpass("Enter database password: ")

# Connect to the database
connection = create_connection("localhost", "root", db_password, "overall_database")

if connection is None:
    print("Failed to connect to the database. Exiting.")
    sys.exit(1)

# Create table (adjust schema as needed)
create_table_query = """
CREATE TABLE IF NOT EXISTS boston_coffee_shops (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ShopName VARCHAR(255),
    Review  FLOAT,
    Ratings VARCHAR(255),
    Address VARCHAR(255),
    Latitude FLOAT,
    Longitude FLOAT
)
"""
execute_query(connection, create_table_query)

# function for inserting variables into table
#
# for some reason GeoJSON stores data in long, lat instead of lat, long. confusing!
def insert_variables_into_table(connection, ShopName, Review, Ratings, Address, Longitude, Latitude):
    try:
        #connection = mysql.connector.connect(host='localhost',database='overall_database',user='root',password=db_password)
        cursor = connection.cursor()
        mySql_insert_query = """INSERT INTO boston_coffee_shops (ShopName, Review, Ratings, Address, Longitude, Latitude)  VALUES (%s, %s, %s, %s, %s, %s) """
        record = (ShopName, Review, Ratings, Address, Longitude, Latitude)
        cursor.execute(mySql_insert_query, record)
        connection.commit()
        print("Record inserted successfully into boston_coffee_shops table")

    except mysql.connector.Error as error:
        print("Failed to insert into MySQL table {}".format(error))

# Insert data
for feature in geojson_data['features']:
    properties = feature['properties']
    coordinates = feature['geometry']
    name = properties["Name"]
    review  = properties["Review"]
    ratings = properties["Ratings"]
    address = properties["Address"]
    longitude = coordinates["coordinates"][0]
    latitude = coordinates["coordinates"][1]
    insert_variables_into_table(connection, name, review, ratings, address, longitude, latitude)
     

print("Data import completed")
