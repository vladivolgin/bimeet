import {
  VStack,
  HStack,
  Text,
  Button,
  Box,
  Avatar,
  AvatarGroup,
  Badge,
  Input,
  NumberInput,
  NumberInputField,
  useDisclosure,
  Collapse,
  FormControl,
  FormLabel,
  Icon,
} from '@chakra-ui/react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { FiMapPin, FiUsers, FiPlus, FiTruck } from 'react-icons/fi'
import type { Event } from '@/entities/event/model/types'
import { useAddCarpool, useJoinCarpool } from '@/features/carpools/model/hooks'
import { useAuth } from '@/features/auth/model/AuthContext'

export default function CarpoolTab({ event }: { event: Event }) {
  const { user } = useAuth()
  const { t } = useTranslation()
  const { isOpen, onToggle } = useDisclosure()
  const [departure, setDeparture] = useState('')
  const [seats, setSeats] = useState('3')

  const addCarpool = useAddCarpool(event.id)
  const joinCarpool = useJoinCarpool(event.id)

  const isDriver = event.carpools.some(c => c.driver.id === user?.id)
  const isCompleted = event.status === 'completed'

  function handleCreate() {
    if (!departure.trim() || !seats) return
    addCarpool.mutate({ departure_point: departure.trim(), seats_available: Number(seats) }, {
      onSuccess: () => {
        setDeparture('')
        setSeats('3')
        onToggle()
      },
    })
  }

  return (
    <VStack align="stretch" spacing={4}>
      <HStack justify="space-between" flexWrap="wrap" gap={2}>
        <Text fontSize="xs" fontWeight="600" color="faintText" textTransform="uppercase" letterSpacing="0.06em">
          {t('carpools.sectionTitle')}
        </Text>
        {!isDriver && !isCompleted && (
          <Button size="sm" leftIcon={<FiPlus />} colorScheme="blue" variant="outline" onClick={onToggle} flexShrink={0}>
            {t('carpools.offerSeats')}
          </Button>
        )}
      </HStack>

      <Collapse in={isOpen} animateOpacity>
        <Box p={4} borderRadius="xl" bg="subtleBg" border="1px solid" borderColor="subtleBorder">
          <VStack spacing={3}>
            <FormControl>
              <FormLabel fontSize="xs" fontWeight="600" color="dimText">{t('carpools.fieldDeparture')}</FormLabel>
              <Input
                size="sm"
                value={departure}
                onChange={e => setDeparture(e.target.value)}
                placeholder={t('carpools.fieldDeparturePlaceholder')}
              />
            </FormControl>
            <FormControl>
              <FormLabel fontSize="xs" fontWeight="600" color="dimText">{t('carpools.fieldSeats')}</FormLabel>
              <NumberInput size="sm" value={seats} onChange={setSeats} min={1} max={8}>
                <NumberInputField />
              </NumberInput>
            </FormControl>
            <HStack w="full" justify="flex-end" spacing={2}>
              <Button size="sm" variant="ghost" onClick={onToggle}>{t('common.cancel')}</Button>
              <Button size="sm" colorScheme="blue" onClick={handleCreate} isLoading={addCarpool.isPending}>{t('common.create')}</Button>
            </HStack>
          </VStack>
        </Box>
      </Collapse>

      {event.carpools.length === 0 && (
        <Box textAlign="center" py={10}>
          <Icon as={FiTruck} boxSize={8} color="faintText" mb={2} />
          <Text fontSize="sm" color="dimText">{t('carpools.empty')}</Text>
        </Box>
      )}

      <VStack align="stretch" spacing={3}>
        {event.carpools.map(carpool => {
          const freeSeats = carpool.seats_available - carpool.passengers.length
          const iAmDriver = carpool.driver.id === user?.id
          const iAmPassenger = carpool.passengers.some(p => p.id === user?.id)

          return (
            <Box
              key={carpool.id}
              p={4}
              borderRadius="xl"
              border="1px solid"
              borderColor="subtleBorder"
              bg="cardBg"
            >
              <HStack justify="space-between" mb={3}>
                <HStack spacing={3} minW={0} flex={1}>
                  <Avatar size="sm" name={carpool.driver.name} bg="brand.400" color="white" flexShrink={0} />
                  <Box minW={0}>
                    <HStack spacing={1.5}>
                      <Text fontWeight="500" fontSize="sm" noOfLines={1}>{carpool.driver.name}</Text>
                      {iAmDriver && <Badge colorScheme="blue" variant="subtle" flexShrink={0}>{t('carpools.badgeYou')}</Badge>}
                    </HStack>
                    <HStack spacing={1} color="dimText">
                      <Icon as={FiMapPin} boxSize={3} />
                      <Text fontSize="xs" noOfLines={1}>{carpool.departure_point}</Text>
                    </HStack>
                  </Box>
                </HStack>
                <HStack spacing={1.5} flexShrink={0}>
                  <Icon as={FiUsers} boxSize={3.5} color="dimText" />
                  <Text
                    fontSize="sm"
                    fontWeight="600"
                    color={freeSeats > 0 ? 'green.600' : 'red.500'}
                  >
                    {freeSeats}/{carpool.seats_available}
                  </Text>
                </HStack>
              </HStack>

              {carpool.passengers.length > 0 && (
                <HStack spacing={2} mb={3}>
                  <Text fontSize="xs" color="dimText">{t('carpools.passengers')}</Text>
                  <AvatarGroup size="xs" max={6} spacing="-6px">
                    {carpool.passengers.map(p => (
                      <Avatar key={p.id} name={p.name} />
                    ))}
                  </AvatarGroup>
                </HStack>
              )}

              {!iAmDriver && !isDriver && !isCompleted && (
                <Button
                  size="sm"
                  w="full"
                  colorScheme={iAmPassenger ? 'red' : 'blue'}
                  variant={iAmPassenger ? 'outline' : 'solid'}
                  isDisabled={!iAmPassenger && freeSeats === 0}
                  onClick={() => user && joinCarpool.mutate(carpool.id)}
                  isLoading={joinCarpool.isPending}
                >
                  {iAmPassenger ? t('carpools.buttonLeave') : freeSeats === 0 ? t('carpools.buttonNoSeats') : t('carpools.buttonJoin')}
                </Button>
              )}
            </Box>
          )
        })}
      </VStack>
    </VStack>
  )
}
