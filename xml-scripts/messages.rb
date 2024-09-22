require 'nori'
require 'pp'

def format_altitude(altitude)
  str = '%.0f' % (altitude / 100)
  str.rjust(3, '0')
end

FACILITY_CODE_MAP = {
  'ZBW' => 'B',
  'ZOB' => 'C',
  'ZTL' => 'T',
  'ZHU' => 'H',
  'ZMA' => 'Z',
  'ZDC'=> 'D'
}

class ATCPosition
  attr_reader :facility, :sector

  def initialize(facility, sector)
    @facility = facility
    @sector = sector
  end

  def format_for_handoff
    facility_code = facility == 'ZJX' ? '-' : (FACILITY_CODE_MAP[facility] || 'Z')
    "#{facility_code}#{sector}"
  end

  def self.from_data(data)
    ATCPosition.new(data['@unitIdentifier'], data['@sectorIdentifier'])
  end
end

class FlightPlan
  # identification fields
  attr_reader :acid, :cid, :identifier
  # departure/arrival
  attr_reader :arrival_airport, :departure_airport
  # tracking facility/sector
  attr_reader :owner
  # handoff sender/receiver (optional)
  attr_reader :handoff_from, :handoff_to
  # position/altitude/speed
  attr_reader :assigned_altitude, :current_altitude, :latitude, :longitude, :speed

  def initialize(flight_data)
    @identifier = flight_data['flightPlan']['@identifier']
    @arrival_airport = flight_data['arrival']['@arrivalPoint']
    @departure_airport = flight_data['departure']['@departurePoint']
    @cid = flight_data['flightIdentification']['@computerId']
    @acid = flight_data['flightIdentification']['@aircraftIdentification']

    if flight_data['controllingUnit']
      facility = flight_data['controllingUnit']['@unitIdentifier']
      sector = flight_data['controllingUnit']['@sectorIdentifier']
      @owner = ATCPosition.from_data(flight_data['controllingUnit'])
    end

    if flight_data['enRoute']
      if flight_data['enRoute']['position']
        @current_altitude = flight_data['enRoute']['position']['altitude'].to_f
        position = flight_data['enRoute']['position']['position']['location']['pos']
        @latitude, @longitude = position.split(' ')
        @speed = flight_data['enRoute']['position']['actualSpeed']['surveillance'].to_f
      end

      if flight_data['enRoute']['boundaryCrossings'] && flight_data['enRoute']['boundaryCrossings']['handoff']
        if flight_data['enRoute']['boundaryCrossings']['handoff']['receivingUnit']
          @handoff_to = ATCPosition.from_data(flight_data['enRoute']['boundaryCrossings']['handoff']['receivingUnit'])
        end
        if flight_data['enRoute']['boundaryCrossings']['handoff']['transferringUnit']
          @handoff_from = ATCPosition.from_data(flight_data['enRoute']['boundaryCrossings']['handoff']['transferringUnit'])
        end
      end
    end

    @assigned_altitude = flight_data['assignedAltitude']['simple'].to_f
  end

  def in_handoff?
    handoff_from || handoff_to
  end

  def to_s
    # Render in a format similar to an ERAM datablock

    # Line 1 (callsign)
    lines = [acid]

    # Line 2 (altitude)
    if current_altitude == assigned_altitude
      lines << "#{format_altitude(current_altitude)}C"
    elsif current_altitude < assigned_altitude
      lines << "#{format_altitude(assigned_altitude)}↑#{format_altitude(current_altitude)}"
    elsif current_altitude > assigned_altitude
      lines << "#{format_altitude(assigned_altitude)}↓#{format_altitude(current_altitude)}"
    end

    # Line 3 (CID/airspeed, or, if in handoff status, CID/handoff)
    if in_handoff?
      if handoff_from
        lines << "#{cid}H#{handoff_from.format_for_handoff}"
      else
        lines << "#{cid}H#{handoff_to.format_for_handoff}"
      end
    else
      lines << "#{cid} #{'%.0f' % speed}"
    end

    # Line 4 (destination airport)
    lines << arrival_airport

    # lines = ["#{acid}: from #{departure_airport} to #{arrival_airport} (cid=#{cid})"]
    #
    # if facility && sector
    #   lines << "Tracked by #{facility} sector #{sector}"
    # end
    lines.join("\n")
  end
end

class FlightPlanManager
  @plans = {} # identifier => FlightPlan
  @last_seen = {} # identifier => DateTime

  def add(plan)
    plans[plan.identifier] = plan
  end
end

def process_flight(flight_data)
  fp = FlightPlan.new(flight_data)
  puts fp
  puts
  # pp flight_data
end

def process_message(message)
  if message['flight']
    process_flight message['flight']
  else
    raise "Got a message without a 'flight' attribute"
  end
end

def process_file(filename)
  doc = Nori.new.parse(File.read(filename))
  messages = doc['ns5:MessageCollection']['message']
  unless messages.is_a? Array
    messages = [messages]
  end
  # process_message messages.first
  messages.each { |message| process_message message }
end

# process_file 'messages/1705188568631.xml'
process_file 'messages/1705188557627.xml'
exit

Dir.glob('messages/*.xml') do |filename|
  puts filename
  process_file filename
  exit
end
