# Pull Latest
git pull

# Build Styles and Hash
npm install -D tailwindcss
npx tailwindcss -i ./public/styles.css -o ./public/styles.min.css --minify

# Generate hash
HASH=$(md5sum public/styles.min.css | cut -d ' ' -f 1)
NEW_CSS_FILE_NAME="styles-$HASH.min.css"

mv ./public/styles.min.css "./public/$NEW_CSS_FILE_NAME"

# Update the .env file with the new hash
if grep -q "BUILD_HASH=" .env; then
    awk -v hash="$HASH" '/BUILD_HASH=/ { $0="BUILD_HASH=" hash } { print }' .env > .env.tmp && mv .env.tmp .env
else
    echo "BUILD_HASH=$HASH" >> .env
fi

echo "Updated build hash to: $HASH"

# Build docker
docker build -t marcellodotgg .
docker kill marcellodotgg_container
docker container prune -f
docker run -d \
       -p 8084:8080 \
       -v /var/www/html:/var/www/html \
       -e GIN_MODE=release \
       --name marcellodotgg_container marcellodotgg
