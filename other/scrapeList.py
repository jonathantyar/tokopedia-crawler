import sys
import time
import logging
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
import re
from datetime import datetime
import json

logging.getLogger('selenium').setLevel(logging.WARNING)

def scrape_page(url):
    # Set Chrome options to disable JavaScript
    chrome_options = Options()
    chrome_options.add_argument("--disable-javascript")

    # Initialize Chrome driver with the options
    driver = webdriver.Chrome(options=chrome_options)

    # Get the URL
    driver.get(url)

    # Wait for the page to load
    time.sleep(5)  # Adjust the sleep time as needed

    # Get the page source
    page_source = driver.page_source

    # Close the driver
    driver.quit()

    # Print or use the page source as needed
    # print(page_source)

    # Extract the value inside window.__memoCache using regex
    pattern = r'window\.__cache\s?=\s?(.*?);(?=\n)'
    match = re.search(pattern, page_source)

    if match:
        # Get the captured group containing the value
        memo_cache_value = match.group(1)

        # Generate the file name with timestamp prefix
        timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
        file_name = f"{timestamp}_cache.json"

        with open(file_name, "w") as file:
            file.write(memo_cache_value)

        json_output = json.dumps({"status": True, "file": file_name, "message": "all is well"})
        print(json_output)
    else:
        json_output = json.dumps({"status": False, "file": "", "message": "window.__cache not found"})
        print(json_output)


if __name__ == "__main__":
    # Check if the 'page' parameter is provided
    if len(sys.argv) > 1:
        category = sys.argv[1]
        page = sys.argv[2]
        # Convert the 'page' parameter to an integer if needed
        try:
            page = int(page)
        except ValueError:
            print("Invalid value for 'page'. Please provide an integer.")
            sys.exit(1)

        # Generate the URL with the desired page number
        url = f"https://www.tokopedia.com/p/{category}/?page={page}"
        scrape_page(url)
    else:
        print("Please provide a value for 'page' parameter.")
        sys.exit(1)