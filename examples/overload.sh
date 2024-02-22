#!/bin/zsh

# Define the splits
hsplits=5
vsplits=5

width=1280
height=720

# Array of colors
random_hex_color() {
    #printf "%02x%02x%02x\n" $((RANDOM % 256)) $((RANDOM % 256)) $((RANDOM % 256))
    echo "$(openssl rand -hex 3)"
}

echo '' > /tmp/colors.txt
echo '' > /tmp/cols.txt
for ((i=1; i<=$hsplits; i++))
do
    echo "$i"
    for ((j=1; j<=$vsplits; j++))
    do
        # Get a random index for selecting color from array
        tmpw=$((width/vsplits))
        tmpx=$((width/vsplits*(j-1)))
        tmph=$((height/hsplits))
        tmpy=$((height/hsplits*(i-1)))
        ebbe color -v -c "$(random_hex_color)" -c 000000 -x $tmpx -y $tmpy -w $tmpw  -h $tmph >> /tmp/colors.txt
    done
done
cat /tmp/colors.txt > /tmp/cols.txt


ebbe text -c ffffff --text acab --fontsize 30 > /tmp/acab.txt
ebbe text --color ffffff -x 112 --text "$ sudo killall -9 fascist" --fontsize  5 >> /tmp/acab.txt
cat /tmp/acab.txt | shuf > /tmp/acab2.txt

ebbe image -y 300 --image enton.png > /tmp/image.png
ebbe image -y 300 -x 200 --image enton.png >> /tmp/image.png
ebbe image -y 300 -x 400 --image enton.png >> /tmp/image.png
ebbe image -y 300 -x 600 --image enton.png >> /tmp/image.png

 {ebbe merge -i /tmp/cols.txt -i /tmp/acab2.txt -i /tmp/image.png; }\
  | ebbe send --oneshot --host pixelflut.uwu.industries:1234 --input - -c 16 -p 1024
