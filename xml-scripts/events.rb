require 'nori'

IGNORED_ATTRIBUTES = [:current_altitude, :latitude, :longitude, :speed]

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

  def to_s
    "#{@facility}#{@sector}"
  end

  def !=(other)
    !other || other == :not_present || facility != other.facility || sector != other.sector
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
  attr_reader :handoff_from, :handoff_to, :handoff_event
  # position/altitude/speed
  attr_reader :assigned_altitude, :current_altitude, :interim_altitude, :latitude, :longitude, :speed
  # flight status
  attr_reader :status
  # fourth line (scratchpad)
  attr_reader :fourth_line_heading, :fourth_line_speed, :fourth_line_text
  # beacon code
  attr_reader :current_beacon_code, :reassigned_beacon_code
  # pointout
  attr_reader :pointout_from, :pointout_to

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
    else
      @owner = $last_fp ? $last_fp.owner : nil
    end

    if flight_data['enRoute']
      if flight_data['enRoute']['position']
        @current_altitude = flight_data['enRoute']['position']['altitude'].to_f
        position = flight_data['enRoute']['position']['position']['location']['pos']
        @latitude, @longitude = position.split(' ')
        @speed = flight_data['enRoute']['position']['actualSpeed']['surveillance'].to_f
      else
        @current_altitude = $last_fp ? $last_fp.current_altitude : nil
        @latitude = $last_fp ? $last_fp.latitude : nil
        @longitude = $last_fp ? $last_fp.longitude : nil
        @speed = $last_fp ? $last_fp.speed : nil
      end

      if flight_data['enRoute']['boundaryCrossings'] && flight_data['enRoute']['boundaryCrossings']['handoff']
        if flight_data['enRoute']['boundaryCrossings']['handoff']['@event']
          @handoff_event = flight_data['enRoute']['boundaryCrossings']['handoff']['@event']
        else
          @handoff_event = $last_fp ? $last_fp.handoff_event : nil
        end
        if flight_data['enRoute']['boundaryCrossings']['handoff']['receivingUnit']
          @handoff_to = ATCPosition.from_data(flight_data['enRoute']['boundaryCrossings']['handoff']['receivingUnit'])
        else
          @handoff_to = $last_fp ? $last_fp.handoff_to : nil
        end
        if flight_data['enRoute']['boundaryCrossings']['handoff']['transferringUnit']
          @handoff_from = ATCPosition.from_data(flight_data['enRoute']['boundaryCrossings']['handoff']['transferringUnit'])
        else
          @handoff_from = $last_fp ? $last_fp.handoff_from : nil
        end
        if handoff_event == "ACCEPTANCE" || handoff_event == "RETRACTION"
          @handoff_from = nil
          @handoff_to = nil
        end
      else
        @handoff_from = $last_fp ? $last_fp.handoff_from : nil
        @handoff_to = $last_fp ? $last_fp.handoff_to : nil
        @handoff_event = $last_fp ? $last_fp.handoff_event : nil
      end

      if flight_data['enRoute']['cleared']
        @fourth_line_heading = flight_data['enRoute']['cleared']['@clearanceHeading']
        @fourth_line_speed = flight_data['enRoute']['cleared']['@clearanceSpeed']
        @fourth_line_text = flight_data['enRoute']['cleared']['@clearanceText']
      else
        @fourth_line_heading = $last_fp ? $last_fp.fourth_line_heading : nil
        @fourth_line_speed = $last_fp ? $last_fp.fourth_line_speed : nil
        @fourth_line_text = $last_fp ? $last_fp.fourth_line_text : nil
      end

      if flight_data['enRoute']['beaconCodeAssignment']
        @current_beacon_code = flight_data['enRoute']['beaconCodeAssignment']['currentBeaconCode']
        @reassigned_beacon_code = flight_data['enRoute']['beaconCodeAssignment']['reassignedBeaconCode']
      else
        @current_beacon_code = $last_fp ? $last_fp.current_beacon_code : nil
        @reassigned_beacon_code = $last_fp ? $last_fp.reassigned_beacon_code : nil
      end

      if flight_data['enRoute']['pointout']
        @pointout_from = ATCPosition.from_data(flight_data['enRoute']['pointout']['originatingUnit'])
        @pointout_to = ATCPosition.from_data(flight_data['enRoute']['pointout']['receivingUnit'])
      else
        @pointout_from = $last_fp ? $last_fp.pointout_from : nil
        @pointout_to = $last_fp ? $last_fp.pointout_to : nil
      end
    else
      @current_altitude = $last_fp ? $last_fp.current_altitude : nil
      @latitude = $last_fp ? $last_fp.latitude : nil
      @longitude = $last_fp ? $last_fp.longitude : nil
      @speed = $last_fp ? $last_fp.speed : nil
      @handoff_from = $last_fp ? $last_fp.handoff_from : nil
      @handoff_to = $last_fp ? $last_fp.handoff_to : nil
      @handoff_event = $last_fp ? $last_fp.handoff_event : nil
      @fourth_line_heading = $last_fp ? $last_fp.fourth_line_heading : nil
      @fourth_line_speed = $last_fp ? $last_fp.fourth_line_speed : nil
      @fourth_line_text = $last_fp ? $last_fp.fourth_line_text : nil
      @current_beacon_code = $last_fp ? $last_fp.current_beacon_code : nil
      @reassigned_beacon_code = $last_fp ? $last_fp.reassigned_beacon_code : nil
      @pointout_from = $last_fp ? $last_fp.pointout_from : nil
      @pointout_to = $last_fp ? $last_fp.pointout_to : nil
    end

    @assigned_altitude = flight_data['assignedAltitude'] ? flight_data['assignedAltitude']['simple'].to_f : :not_present

    if flight_data['interimAltitude']
      @interim_altitude = flight_data['interimAltitude'].is_a?(String) ? flight_data['interimAltitude'] : nil
    else
      @interim_altitude = $last_fp ? $last_fp.interim_altitude : nil
    end

    if flight_data['flightStatus']
      @status = flight_data['flightStatus']['@fdpsFlightStatus']
    else
      @status = $last_fp ? $last_fp.status : nil
    end
  end

  def is_changed(from, to)
    from != to && from != :not_present && to != :not_present
  end

  def diff(fp)
    changes = {}
    if is_changed(acid, fp.acid)
      changes[:acid] = [acid, fp.acid]
    end
    if is_changed(cid, fp.cid)
      changes[:cid] = [cid, fp.cid]
    end
    if is_changed(identifier, fp.identifier)
      changes[:identifier] = [identifier, fp.identifier]
    end
    if is_changed(arrival_airport, fp.arrival_airport)
      changes[:arrival_airport] = [arrival_airport, fp.arrival_airport]
    end
    if is_changed(departure_airport, fp.departure_airport)
      changes[:departure_airport] = [departure_airport, fp.departure_airport]
    end
    if is_changed(owner, fp.owner)
      changes[:owner] = [owner, fp.owner]
    end
    if is_changed(handoff_from, fp.handoff_from)
      changes[:handoff_from] = [handoff_from, fp.handoff_from]
    end
    if is_changed(handoff_to, fp.handoff_to)
      changes[:handoff_to] = [handoff_to, fp.handoff_to]
    end
    if is_changed(handoff_event, fp.handoff_event)
      changes[:handoff_event] = [handoff_event, fp.handoff_event]
    end
    if is_changed(assigned_altitude, fp.assigned_altitude)
      changes[:assigned_altitude] = [assigned_altitude, fp.assigned_altitude]
    end
    if is_changed(current_altitude, fp.current_altitude)
      changes[:current_altitude] = [current_altitude, fp.current_altitude]
    end
    if is_changed(interim_altitude, fp.interim_altitude)
      changes[:interim_altitude] = [interim_altitude, fp.interim_altitude]
    end
    if is_changed(status, fp.status)
      changes[:status] = [status, fp.status]
    end
    if is_changed(latitude, fp.latitude)
      changes[:latitude] = [latitude, fp.latitude]
    end
    if is_changed(longitude, fp.longitude)
      changes[:longitude] = [longitude, fp.longitude]
    end
    if is_changed(speed, fp.speed)
      changes[:speed] = [speed, fp.speed]
    end
    if is_changed(fourth_line_heading, fp.fourth_line_heading)
      changes[:fourth_line_heading] = [fourth_line_heading, fp.fourth_line_heading]
    end
    if is_changed(fourth_line_speed, fp.fourth_line_speed)
      changes[:fourth_line_speed] = [fourth_line_speed, fp.fourth_line_speed]
    end
    if is_changed(fourth_line_text, fp.fourth_line_text)
      changes[:fourth_line_text] = [fourth_line_text, fp.fourth_line_text]
    end
    if is_changed(current_beacon_code, fp.current_beacon_code)
      changes[:current_beacon_code] = [current_beacon_code, fp.current_beacon_code]
    end
    if is_changed(reassigned_beacon_code, fp.reassigned_beacon_code)
      changes[:reassigned_beacon_code] = [reassigned_beacon_code, fp.reassigned_beacon_code]
    end
    if pointout_from != fp.pointout_from
      changes[:pointout_from] = [pointout_from, fp.pointout_from]
    end
    if pointout_to != fp.pointout_to
      changes[:pointout_to] = [pointout_to, fp.pointout_to]
    end

    IGNORED_ATTRIBUTES.each do |key|
      changes.delete(key)
    end

    changes
  end
end

def value_or_nil(value)
  value ? value : '<nil>'
end

def show_changes(changes)
  return if changes.empty?
  changes.each do |key, value|
    puts "#{key} changed from #{value_or_nil(value[0])} to #{value_or_nil(value[1])}"
  end
end

def process_flight(flight_data, filename)
  new_fp = FlightPlan.new(flight_data)
  if $last_fp
    changes = $last_fp.diff(new_fp)
    unless changes.empty?
      puts "--------"
      puts filename
      show_changes(changes)
    end
  end
  $last_fp = new_fp
end

def process_message(filename, message, callsign)
  if message['flight'] && message['flight']['flightIdentification'] && message['flight']['flightIdentification']['@aircraftIdentification'] == callsign
    process_flight(message['flight'], filename)
  end
end

def process_file(filename, callsign)
  doc = Nori.new.parse(File.read(filename))
  return unless doc['ns5:MessageCollection']
  messages = doc['ns5:MessageCollection']['message']
  unless messages.is_a? Array
    messages = [messages]
  end
  messages.each { |message| process_message(filename, message, callsign) }
end

if ARGV.length != 1
  $stderr.puts 'USAGE: events.rb <callsign>'
  exit
end

callsign = ARGV.first

Dir.glob('messages/*.xml').sort.each do |file|
  begin
    process_file(file, callsign)
  rescue => e
    $stderr.puts "Error processing file #{file}"
    raise e
  end
end

