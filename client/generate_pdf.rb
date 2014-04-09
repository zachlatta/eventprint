#!/usr/bin/env ruby
#
# This script expects the following unnamed arguments in their respective
# order:
#
# barcode, first name, last name, title, and file name

require 'prawn'
require 'prawn/qrcode'

BARCODE = ARGV[0]
FIRST_NAME = ARGV[1]
LAST_NAME = ARGV[2]
TITLE = ARGV[3]
HACKATHON_LOGO_PATH = 'lahacks.png'
FILE_NAME = ARGV[4]
 
Prawn::Document.generate(FILE_NAME, page_size: [200, 350]) do
  font 'Helvetica'
  pad_top(112) do
    image HACKATHON_LOGO_PATH, width: 90, at: [20,275]
    text FIRST_NAME, style: :bold, size: 24, align: :center, top_margin: '100'
    text LAST_NAME, style: :bold, size: 24, align: :center
  end
  pad_top(2) do
    print_qr_code(BARCODE, extent: 100, align: :center)
  end
  fill { rectangle [-36,10], 200, 50 }
  bounding_box([-36,-3], width: 200, height: 40) do
    text TITLE, style: :bold, size: 20, color: 'ffffff', align: :center
  end
end
