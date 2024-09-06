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
with open('data.json') as file:
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
    name VARCHAR(255),
    review  FLOAT,
    ratings INT,
    address VARCHAR(255),
    latitude FLOAT,
    longitude FLOAT
)
"""
execute_query(connection, create_table_query)

# Insert data
for feature in geojson_data['features']:
    properties = json.dumps(feature['properties'])
    geometry = json.dumps(feature['geometry'])

    name = properties['key']
    review  = properties['Review']
    ratings = properties['Ratings']
    address = properties['Address']
    longitude, latitude = geometry['coordinates']

    insert_query = """
    INSERT INTO boston_coffee_shops (name, review, ratings, address, latitude, longitude)
    VALUES (%s, %s, %s, %s, %s, %s) 
    """
    execute_query(connection, insert_query)

print("Data import completed")
