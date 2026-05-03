import {
  Box,
  Button,
  Heading,
  HStack,
  Text,
  VStack,
  Icon,
  Avatar,
  AvatarGroup,
  Badge,
  Tooltip,
} from '@chakra-ui/react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { FiCalendar, FiGlobe, FiLock, FiMapPin, FiUsers } from 'react-icons/fi'
import type { EventListItem } from '@/entities/event/model/types'
import { formatDate, formatTime } from '@/shared/lib/formatDate'

const CARD_GRADIENTS = [
  'linear(135deg, #4F46E5, #7C3AED)',
  'linear(135deg, #F43F5E, #A855F7)',
  'linear(135deg, #0891B2, #4F46E5)',
  'linear(135deg, #059669, #0891B2)',
  'linear(135deg, #D97706, #E11D48)',
]

const BUSINESS_GRADIENT = 'linear(135deg, #0F172A, #1E3A8A)'

export function getGradient(id: string, category?: string) {
  if (category === 'business') return BUSINESS_GRADIENT
  const idx = id.charCodeAt(id.length - 1) % CARD_GRADIENTS.length
  return CARD_GRADIENTS[idx]
}

const statusScheme = {
  invited:   'yellow',
  confirmed: 'green',
  declined:  'red',
} as const

interface EventCardProps {
  event: EventListItem
  myStatus?: string
  canJoin?: boolean
  onJoin?: (e: React.MouseEvent) => void
}

export default function EventCard({ event, myStatus, canJoin, onJoin }: EventCardProps) {
  const { t } = useTranslation()
  const confirmedCount = event.confirmed_count ?? event.participants.filter(p => p.status === 'confirmed').length
  const isCompleted = event.status === 'completed'
  const gradient = getGradient(event.id, event.category)
  const statusKey = myStatus as keyof typeof statusScheme | undefined

  return (
    <Box
      as={Link}
      to={`/events/${event.id}`}
      display="flex"
      flexDirection="column"
      bg="cardBg"
      borderRadius="xl"
      overflow="hidden"
      border="1px solid"
      borderColor="cardBorder"
      boxShadow="0 1px 3px rgba(15,23,42,0.04), 0 4px 16px rgba(15,23,42,0.05)"
      _hover={{
        transform: 'translateY(-2px)',
        boxShadow: '0 2px 8px rgba(15,23,42,0.06), 0 12px 32px rgba(15,23,42,0.09)',
        borderColor: 'defaultBorder',
      }}
      transition="all 0.2s cubic-bezier(0.4, 0, 0.2, 1)"
    >
      {/* Accent bar — gray when completed */}
      <Box
        h="4px"
        flexShrink={0}
        {...(isCompleted
          ? { bg: 'gray.300', _dark: { bg: 'gray.600' } }
          : { bgGradient: gradient }
        )}
      />

      {/* Card body — dimmed when completed */}
      <Box p={5} flex={1} display="flex" flexDirection="column" opacity={isCompleted ? 0.6 : 1}>

        {/* Header: badges + title + status */}
        <Box mb={3}>
          <HStack justify="space-between" align="flex-start">
            <Box flex={1} mr={2}>
              <HStack mb={1} spacing={1.5} flexWrap="wrap">
                {/* Completed badge — first and solid */}
                {isCompleted && (
                  <Badge colorScheme="gray" variant="solid" fontSize="9px" borderRadius="md">
                    {t('eventCard.statusCompleted')}
                  </Badge>
                )}
                <Badge
                  colorScheme={event.category === 'business' ? 'gray' : 'green'}
                  variant="subtle"
                  fontSize="9px"
                  borderRadius="md"
                >
                  {event.category === 'business' ? t('eventCard.categoryBusiness') : t('eventCard.categoryOrdinary')}
                </Badge>
                <Badge
                  colorScheme={event.is_public ? 'teal' : 'purple'}
                  variant="subtle"
                  fontSize="9px"
                  borderRadius="md"
                  display="flex"
                  alignItems="center"
                  gap={0.5}
                >
                  <Icon as={event.is_public ? FiGlobe : FiLock} boxSize={2.5} />
                  {event.is_public ? t('eventCard.publicBadge') : t('eventCard.personalBadge')}
                </Badge>
              </HStack>
              <Heading
                size="sm"
                noOfLines={2}
                letterSpacing="-0.2px"
                lineHeight="1.3"
              >
                {event.title}
              </Heading>
            </Box>
            {statusKey && (
              <Badge
                colorScheme={statusScheme[statusKey]}
                variant="subtle"
                flexShrink={0}
                alignSelf="flex-start"
              >
                {t(`eventCard.status${statusKey.charAt(0).toUpperCase()}${statusKey.slice(1)}`)}
              </Badge>
            )}
          </HStack>
        </Box>

        {/* Meta — date + location */}
        <VStack align="stretch" spacing={1.5}>
          <HStack spacing={2}>
            <Icon as={FiCalendar} color="brand.400" boxSize={3.5} flexShrink={0} />
            <Text fontSize="xs" color="dimText" fontWeight="500" noOfLines={1}>
              {formatDate(event.date_start)}, {formatTime(event.date_start)}
            </Text>
          </HStack>
          <HStack spacing={2}>
            <Icon as={FiMapPin} color="brand.400" boxSize={3.5} flexShrink={0} />
            <Text fontSize="xs" color="dimText" fontWeight="500" noOfLines={1}>
              {event.location}
            </Text>
          </HStack>
        </VStack>

        {/* Spacer pushes footer to bottom */}
        <Box flex={1} />

        {/* Footer */}
        <HStack justify="space-between" mt={3} pt={3} borderTop="1px solid" borderColor="subtleBorder">
          <HStack spacing={1.5}>
            <Icon as={FiUsers} color="faintText" boxSize={3.5} />
            <Text fontSize="xs" color="faintText">{confirmedCount}</Text>
            <AvatarGroup size="xs" max={3} spacing="-6px">
              {event.participants
                .filter(p => p.status === 'confirmed')
                .map(p => (
                  <Avatar key={p.id} name={p.user.name} bg="brand.400" color="white" />
                ))}
            </AvatarGroup>
          </HStack>
          {canJoin != null ? (
            <Tooltip label={!canJoin ? t('eventInfo.noSpots') : undefined} isDisabled={canJoin}>
              <Button
                size="xs"
                colorScheme="blue"
                isDisabled={!canJoin}
                onClick={(e) => { e.preventDefault(); e.stopPropagation(); onJoin?.(e) }}
              >
                {t('eventCard.attendButton')}
              </Button>
            </Tooltip>
          ) : (
            <Text fontSize="xs" fontWeight="500" color="faintText">→</Text>
          )}
        </HStack>
      </Box>
    </Box>
  )
}
