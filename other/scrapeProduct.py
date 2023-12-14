import sys
import json
import time
from selenium import webdriver
from selenium.webdriver.common.by import By

# Check if the URL argument is provided
if len(sys.argv) < 2:
    print("Usage: python scraper.py <url>")
    sys.exit(1)

# Get the URL from the command line argument
url = sys.argv[1]

# Set up Selenium WebDriver (make sure you have the appropriate WebDriver for your browser installed)
driver = webdriver.Chrome()

# Navigate to the page
driver.get(url)

# Wait for the page to load
time.sleep(5)  # Adjust the sleep time as needed

div_element_desc = driver.find_element(By.CSS_SELECTOR, 'div[data-testid="lblPDPDescriptionProduk"]')
description = div_element_desc.get_attribute('innerText')

div_element_name = driver.find_element(By.CSS_SELECTOR, 'h1[data-testid="lblPDPDetailProductName"]')
name = div_element_name.get_attribute('innerText')

div_element_price = driver.find_element(By.CSS_SELECTOR, 'div[data-testid="lblPDPDetailProductPrice"]')
price = div_element_price.get_attribute('innerText')

div_element_rating = driver.find_element(By.CSS_SELECTOR, 'span[data-testid="lblPDPDetailProductRatingNumber"]')
rating = div_element_rating.get_attribute('innerText')

div_element_img = driver.find_element(By.CSS_SELECTOR, 'img[data-testid="PDPMainImage"]')
img = div_element_img.get_attribute('src')

div_element_merchant = driver.find_element(By.CSS_SELECTOR, 'h2[class="css-1wdzqxj-unf-heading e1qvo2ff2"]')
merchant = div_element_merchant.get_attribute('innerText')
# Wrap the content in a JSON object
output = json.dumps({'name':name,'description': description, 'prices': price, 'rating': float(rating), 'image_link': img, 'merchant': merchant})

# Print the JSON string
print(output)

# Close the WebDriver
driver.quit()