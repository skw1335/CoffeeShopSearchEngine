from collections import defaultdict
from webdriver_manager.chrome import ChromeDriverManager
from selenium.webdriver.common.by import By
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.common.keys import Keys
from selenium import webdriver
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import time
import csv

# Initalize a bunch of arrays to store each individual attribute

Names = []
Reviews = []
Stars = []
Stars_count = []
Locations = []
Addresses = []

database = defaultdict(list)

url_list = [ "https://www.google.com/maps/search/coffee+shop+in+boston/",  
   "https://www.google.com/maps/search/coffee+shop+in+brookline/",  
   "https://www.google.com/maps/search/coffee+shop+in+cambridge/" ]

Overall_database = {}
# Initialize the webdriver
driver = webdriver.Chrome()
def WebScraper(url):
    Names = []
    Reviews = []
    Stars = []
    Stars_count = []
    Locations = []
    Addresses = []


    # Open the desired webpage
    driver.get(url)
# Wait for the main elements to be present
    WebDriverWait(driver, 10).until(EC.presence_of_all_elements_located((By.CLASS_NAME, "hfpxzc")))

# Scroll to load more elements

# Now find the elements
    actions = ActionChains(driver)
    for _ in range(200):
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
    counter = 0
    for i in range(min(len(Stars), len(Names))):
       # print(f" index = {i} | Names[i] = {Names[i]}")
       # print(f" index = {i} | Stars[i] = {Stars[i]}")
       # print(f" index = {i} | Addresses[i] = {Addresses[i]}")
        database[Names[i]].append([Stars[i], Addresses[i]])
 



   
    print("we made it")
    #clear the lists used to store everything at the end of every run through
    
    print("to the end of the function :)") 
def WebScraper_iter():
    for url in url_list:
        WebScraper(url)
    with open('Overall_database.csv', 'w') as f:
        w = csv.writer(f)
        for key, value in database.items():
                w.writerow([key, value])
                w.writerow("\n")




WebScraper_iter()
