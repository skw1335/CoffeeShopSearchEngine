from collections import defaultdict
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.common.by import By
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.common.keys import Keys
from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from shapely.geometry import Point
from shapely_geojson import dumps, Feature, FeatureCollection
import time
import csv
import pymysql
import mysql.connector
import json




database = defaultdict(list)
## this can be changed to anything after /maps/search/
url_list = [ "https://www.google.com/maps/search/coffee+shop+in+boston/",  
   "https://www.google.com/maps/search/coffee+shop+in+brookline/",  
   "https://www.google.com/maps/search/coffee+shop+in+cambridge/",
   ]

Overall_database = {}
# Initialize the webdriver
driver = webdriver.Chrome()

#function to add database rows to mysql database
#function honestly doesn't need to exist, saving it just incase
def insert_variables_into_table(ShopName, lat, lng, Reviews, Ratings, Address):
    try:
        connection = mysql.connector.connect(host='localhost',database='overall_database',user='sam',password='')
        cursor = connection.cursor()
        mySql_insert_query = """INSERT INTO boston (ShopName, lat, lng, Reviews, Ratings, Address) VALUES (%s, %s, %s, %s, %s, %s) """
        record = (ShopName, lat, lng, Reviews, Ratings, Address)
        cursor.execute(mySql_insert_query, record)
        connection.commit()
        print("Record inserted successfully into boston table")

    except mysql.connector.Error as error:
        print("Failed to insert into MySQL table {}".format(error))

        
def WebScraper(url):
    Names = []
    Reviews = []
    Stars = []
    Stars_count = []
    Coordinates = []
    Locations = []
    Addresses = []
    output_list = []
    # Open the desired webpage
    driver.get(url)
# Wait for the main elements to be present
    WebDriverWait(driver, 10).until(EC.presence_of_all_elements_located((By.CLASS_NAME, "hfpxzc")))

# Scroll to load more elements

# Now find the elements
    actions = ActionChains(driver)
    for _ in range(150):
        names = driver.find_elements(By.CLASS_NAME, "hfpxzc")
        reviews_and_ratings = driver.find_elements(By.CLASS_NAME, "ZkP5Je")
        info = driver.find_elements(By.CLASS_NAME, "W4Efsd")
        actions.send_keys(3*Keys.ARROW_DOWN).perform()
        time.sleep(1)
# Get the different data parameters, Name, Rating and reviews, Address
    for index, name in enumerate(names): 
        store_name = name.get_attribute("aria-label")
        Names.append(store_name)
    for index, r_and_r in enumerate(reviews_and_ratings):
        review_and_rating = r_and_r.get_attribute("aria-label")
        Stars.append(review_and_rating)
    for i, parent in enumerate(info):
        child = parent.find_element(By.XPATH, "./*")
        if child.tag_name == 'div':
            location = child.get_attribute("innerText")
            location = location.split(',')
            Locations.append(location)
    for index, elem in enumerate(Locations):
        if index % 2 != 0:
            s = ''.join(elem)
            tmp = s.split('·')
            Addresses.append(tmp[-1])
    for i in range(min(len(Stars), len(Names))):
       # print(f" index = {i} | Names[i] = {Names[i]}")
       # print(f" index = {i} | Stars[i] = {Stars[i]}")
       # print(f" index = {i} | Addresses[i] = {Addresses[i]}")
       # print(f" index = {i} | Coordinates[i] = {Coordinates[i]}")
        s = Stars[i].split(' ')
        database[Names[i]].append([Stars[i], Addresses[i]])
    

    print("we made it")
    
    print("to the end of the function :)")
def WebScraper_iter():
    for url in url_list:
        WebScraper(url)

WebScraper_iter()
tmp = []
features = []
json_data_with_duplicates = {}
json_data = {}
def Coord_Scraper():
    ind = 0
    for key, value in database.items():
        for val in value:
            driver.get("https://google.com/maps/search/" + key + ' ' + val[1].strip())
            coordinates = (driver.find_elements(By.CSS_SELECTOR, 'meta[itemprop=image]'))
            for i in coordinates:
                x = i.get_attribute('content') 
                tmp.append(x)
            for i in tmp:
                i = i.split('?center=')[1].split('&zoom=')[0].split('%2C')
            
            s = val[0].split(' ')
            #print(f" key = {key} | i[0] = {i[0]} | i[1] = {i[1]} | s[0] = {s[0]} | s[2] = {s[2]} | val[1] = {val[1]} ")
            feature = Feature(Point(i[1], i[0]), 
                    { 'Name': key,
                      'Review': s[0],
                      'Ratings': s[2],
                      'Address': val[1]
                    })
            features.append(feature)
            ind += 1

    feature_collection = FeatureCollection(features)
    print(dumps(feature_collection, indent=2))
    with open('data.json', 'w') as f:
        f.write(dumps(feature_collection, indent=2))
Coord_Scraper()


