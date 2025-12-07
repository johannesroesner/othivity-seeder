# Othivity Seeder

*Othivity Seeder* is a tool designed to populate our  [Othivity WebApp](https://github.com/johannesroesner/othivity) with data for demonstration purposes. It generates and inserts sample data into the application's database, allowing users to explore and interact with the app's features without needing to input their own data.

## How To Use Othivity Seeder

1. **Clone the Repository**: Start by cloning the *Othivity Seeder* repository to your local machine.

   ```bash
   git clone https://github.com/johannesroesner/othivity-seeder.git
   
2. **Navigate to the Directory**: Change into the cloned directory.

   ```bash
   cd othivity-seeder
   ```
   
3. **Start the seeder**: Ensure that your Othivity WebApp is running locally, and you have a valid jwt generated for a MODERATOR profile. Then, execute the seeder script to populate the database with sample data.

   ```bash
   go run main.go <your jwt token> <target base url>
   ```

## Data and Credits
The address data used by *Othivity Seeder* is sourced from the [Amtliches Straßenverzeichnis Regensburg](https://www.regensburg.de/fm/121/amtl_verzeichnisse_stadtbezirke_nach_strassen.pdf) project. Zipcodes where gathered via [OpenPLZ API](https://www.openplzapi.org/de/). The other sample data e.g. names, activity titles, and descriptions are generated with AI.