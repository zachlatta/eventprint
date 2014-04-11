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

def shorten (string)
  if string.length > 10
    return string[0..9]
  else
    return string
  end
end

Prawn::Document.generate(FILE_NAME, page_size: [144, 252]) do
  font 'Helvetica'
  pad_top(75) do
    image HACKATHON_LOGO_PATH, width: 70, at: [0,190]
    text shorten(FIRST_NAME), style: :bold, size: 12, align: :center, top_margin: '40'
    text shorten(LAST_NAME), style: :bold, size: 12, align: :center
  end
  pad_top(5) do
    print_qr_code(BARCODE, extent: 72, align: :center)
  end
  fill { rectangle [-36,-6], 144, 30 }
  bounding_box([-36,-13], width: 144, height: 40) do
    text TITLE, style: :bold, size: 16, color: 'ffffff', align: :center
  end
end
