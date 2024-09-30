import json
import psycopg2
from psycopg2 import Error
import getpass
import sys

# Function to create database connection
def create_connection(host_name, user_name, user_password, db_name):
    connection = None
    try:
        connection = psycopg2.connect(
            host=host_name,
            user=user_name,
            password=user_password,
            database=db_name
        )
        print("Connection to PostgreSQL DB successful")
    except (Exception, Error) as e:
        print(f"The error '{e}' occurred")
    return connection

# Function to execute query
def execute_query(connection, query):
    cursor = connection.cursor()
    try:
        cursor.execute(query)
        connection.commit()
        print("Query executed successfully")
    except (Exception, Error) as e:
        print(f"The error '{e}' occurred")

# Read GeoJSON file
with open('data_harvest/deduped_data.json') as file:
    geojson_data = json.load(file)

# Take password as CLI
db_password = getpass.getpass("Enter database password: ")

# Connect to the database
connection = create_connection("localhost", "postgres", db_password, "coffeeMap")
if connection is None:
    print("Failed to connect to the database. Exiting.")
    sys.exit(1)

# Create table (adjust schema as needed)
create_table_query = """
CREATE TABLE IF NOT EXISTS coffee_shops (
    id SERIAL PRIMARY KEY,
    shop_name VARCHAR(255),
    review FLOAT,
    ratings VARCHAR(255),
    address VARCHAR(255),
    latitude FLOAT,
    longitude FLOAT
)
"""
execute_query(connection, create_table_query)

# Function for inserting variables into table
def insert_variables_into_table(connection, shop_name, review, ratings, address, longitude, latitude):
    try:
        cursor = connection.cursor()
        postgres_insert_query = """ 
        INSERT INTO coffee_shops (shop_name, review, ratings, address, longitude, latitude)
        VALUES (%s, %s, %s, %s, %s, %s) 
        """
        record = (shop_name, review, ratings, address, longitude, latitude)
        cursor.execute(postgres_insert_query, record)
        connection.commit()
        print("Record inserted successfully into coffee_shops table")
    except (Exception, Error) as error:
        print("Failed to insert into PostgreSQL table: {}".format(error))

# Insert data
for feature in geojson_data['features']:
    properties = feature['properties']
    coordinates = feature['geometry']
    name = properties["Name"]
    review = properties["Review"]
    ratings = properties["Ratings"]
    address = properties["Address"]
    longitude = coordinates["coordinates"][0]
    latitude = coordinates["coordinates"][1]
    insert_variables_into_table(connection, name, review, ratings, address, longitude, latitude)

print("Data import completed")
