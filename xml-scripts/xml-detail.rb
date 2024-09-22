require 'nori'
require 'pp'

puts "Ruby script has started..."

EXPECTED_ATTRS = {
  root: {
    'flight' => {
      '@centre' => true,
      '@flightType' => true,
      '@source' => true,
      '@system' => true,
      '@timestamp' => true,
      'aircraftDescription' => {
        '@aircraftAddress' => true,
        '@aircraftPerformance' => true,
        '@aircraftQuantity' => true,
        '@equipmentQualifier' => true,
        '@registration' => true,
        '@tfmsSpecialAircraftQualifier' => true,
        '@wakeTurbulence' => true,
        'accuracy' => ['cmsFieldType'],
        'aircraftType' => ['icaoModelIdentifier', 'otherModelData'],
        'capabilities' => {
          '@standardCapabilities' => true,
          'communication' => ['@otherCommunicationCapabilities', '@otherDataLinkCapabilities', '@selectiveCallingCode', 'communicationCode', 'dataLinkCode'],
          'navigation' => ['@otherNavigationCapabilities', 'navigationCode', 'performanceBasedCode'],
          'surveillance' => ['@otherSurveillanceCapabilities', 'surveillanceCode'],
        },
      },
      'agreed' => {
        'route' => {
          '@atcIntendedRoute' => true,
          '@flightDuration' => true,
          '@initialFlightRules' => true,
          '@localIntendedRoute' => true,
          '@nasRouteText' => true,
          'adaptedArrivalDepartureRoute' => ['@nasRouteAlphanumeric', '@nasRouteIdentifier'],
          'adaptedDepartureRoute' => ['@nasRouteAlphanumeric', '@nasRouteIdentifier'],
          'estimatedElapsedTime' => {
            '@elapsedTime' => true,
            'location' => {
              'point' => ['@fix'],
              'region' => true,
            },
          },
          'expandedRoute' => {
            'routePoint' => true,
          },
          'holdFix' => ['@fix', 'distance', 'radial'],
          'nasadaptedArrivalRoute' => ['@nasRouteAlphanumeric', '@nasRouteIdentifier', 'nasFavNumber'],
        },
      },
      'arrival' => {
        '@arrivalPoint' => true,
        'arrivalAerodrome' => {
          '@code' => true,
          '@name' => true,
          'point' => {
            'location' => ['@srsName', 'pos'],
          },
        },
        'arrivalAerodromeAlternate' => ['@code', '@name'],
        'runwayPositionAndTime' => {
          'runwayTime' => {
            'actual' => ['@time'],
            'estimated' => ['@time'],
          },
        },
      },
      'assignedAltitude' => {
        'block' => {
          'above' => ['@uom'],
          'below' => ['@uom'],
        },
        'simple' => true,
        'vfr' => true,
        'vfrOnTopPlus' => true,
        'vfrPlus' => true,
      },
      'controllingUnit' => ['@sectorIdentifier', '@unitIdentifier'],
      'coordination' => {
        '@coordinationTime' => true,
        '@coordinationTimeHandling' => true,
        '@delayTimeToAbsorb' => true,
        'coordinationFix' => {
          '@fix' => true,
          'distance' => true,
          'location' => {
            '@srsName' => true,
            'pos' => true,
          },
          'radial' => true,
        },
      },
      'departure' => {
        '@departurePoint' => true,
        'departureAerodrome' => {
          'point' => {
            'location' => ['@srsName', 'pos'],
          },
        },
        'runwayPositionAndTime' => {
          'runwayTime' => {
            'actual' => ['@time'],
            'controlled' => ['@time'],
            'estimated' => ['@time'],
          },
        },
        'takeoffAlternateAerodrome' => ['@code'],
      },
      'enRoute' => {
        'alternateAerodrome' => ['@code'],
        'beaconCodeAssignment' => ['currentBeaconCode', 'previousBeaconCode', 'reassignedBeaconCode'],
        'boundaryCrossings' => {
          'handoff' => {
            '@event' => true,
            'acceptingUnit' => ['@sectorIdentifier', '@unitIdentifier'],
            'receivingUnit' => ['@sectorIdentifier', '@unitIdentifier'],
            'transferringUnit' => ['@sectorIdentifier', '@unitIdentifier'],
          },
        },
        'cleared' => ['@clearanceHeading', '@clearanceSpeed', '@clearanceText'],
        'expectedFurtherClearanceTime' => ['@time'],
        'pointout' => {
          'originatingUnit' => ['@sectorIdentifier', '@unitIdentifier'],
          'receivingUnit' => ['@sectorIdentifier', '@unitIdentifier'],
        },
        'position' => {
          '@coastIndicator' => true,
          '@positionTime' => true,
          '@reportSource' => true,
          '@targetPositionTime' => true,
          'actualSpeed' => ['surveillance'],
          'altitude' => true,
          'position' => {
            'location' => { '@srsName' => true, 'pos' => true },
          },
          'targetAltitude' => true,
          'targetPosition' => ['@srsName', 'pos'],
          'trackVelocity' => ['x', 'y'],
        },
      },
      'flightIdentification' => ['@aircraftIdentification', '@computerId', '@siteSpecificPlanId'],
      'flightIdentificationPrevious' => ['@aircraftIdentification', '@computerId', '@siteSpecificPlanId'],
      'flightPlan' => ['@flightPlanRemarks', '@identifier'],
      'flightStatus' => ['@airborneHold', '@fdpsFlightStatus'],
      'interimAltitude' => ['@uom'],
      'gufi' => true,
      'operator' => {
        'operatingOrganization' => {
          'organization' => ['@name'],
        },
      },
      'originator' => ['aftnAddress', 'flightOriginator'],
      'requestedAirspeed' => ['nasAirspeed'],
      'requestedAltitude' => ['simple', 'vfr', 'vfrPlus'],
      'routeToRevisedDestination' => {
        'route' => ['@routeText'],
      },
      'specialHandling' => true,
      'supplementalData' => {
        'additionalFlightInformation' => ['nameValue'],
      },
    },
  }
}

def process_node(message, expected_attrs, chain = [])
  here = chain.join('.')
  return unless message

  message.keys.each do |key|
    if key.start_with?('@xsi') || key.start_with?('@xmlns')
      next
    end
    if expected_attrs.key?(key)
      if message[key].is_a?(Hash)
        if expected_attrs[key].is_a?(Hash)
          process_node message[key], expected_attrs[key], chain + [key]
        elsif expected_attrs[key].is_a?(Array)
          keys = expected_attrs[key].map { |key| [key, true] }.to_h
          process_node message[key], keys, chain + [key]
        else
          puts "Found a #{here}.#{key} attribute with unexpected children: #{message[key].keys.sort.join(', ')}"
        end
      else
        puts "Attribute #{here}.#{key} matches expected structure."
      end
    else
      puts "Unexpected attribute: #{here}.#{key}"
    end
  end
end

def process_file(filename)
  puts "Reading file: #{filename}"
  content = File.read(filename)
  puts "File content loaded, length: #{content.length}"

  doc = Nori.new.parse(content)
  puts "XML parsed successfully!"

  collection = doc['ns5:MessageCollection']
  if collection
    puts "MessageCollection found"
    messages = collection['message']
    if messages.is_a?(Array)
      puts "Processing #{messages.size} messages"
      messages.each do |message|
        process_node(message, EXPECTED_ATTRS[:root])
      end
    else
      puts "Single message found"
      process_node(messages, EXPECTED_ATTRS[:root])
    end
  else
    puts "No MessageCollection found in the XML!"
  end
end

if ARGV.length != 1
  puts "Usage: ruby xml-detail.rb <filename>"
else
  process_file(ARGV[0])
end