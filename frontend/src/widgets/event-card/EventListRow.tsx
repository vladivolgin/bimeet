import {
  Badge,
  Box,
  Button,
  HStack,
  Text,
  Icon,
  Avatar,
  AvatarGroup,
  Tooltip,
} from '@chakra-ui/react'
import { Link } from 'react-router-dom'
import { useTranslation } from 'react-i18next'
import { FiCalendar, FiChevronRight, FiGlobe, FiLock, FiMapPin } from 'react-icons/fi'
import type { EventListItem } from '@/entities/event/model/types'
import { formatDate, formatTime } from '@/shared/lib/formatDate'
import { getGradient } from './index'

const statusColors = {
  invited:   { color: '#F59E0B', bg: '#FFFBEB' },
  confirmed: { color: '#10B981', bg: '#ECFDF5' },
  declined:  { color: '#EF4444', bg: '#FEF2F2' },
}

interface EventListRowProps {
  event: EventListItem
  myStatus?: string
  canJoin?: boolean
  onJoin?: (e: React.MouseEvent) => void
}

export default function EventListRow({ event, myStatus, canJoin, onJoin }: EventListRowProps) {
  const { t } = useTranslation()
  const confirmedCount = event.confirmed_count ?? event.participants.filter(p => p.status === 'confirmed').length
  const isCompleted = event.status === 'completed'
  const gradient = getGradient(event.id, event.category)
  const statusKey = myStatus as keyof typeof statusColors | undefined
  const colors = statusKey ? statusColors[statusKey] : null

  return (
    <Box
      as={Link}
      to={`/events/${event.id}`}
      display="grid"
      gridTemplateColumns={{
        base: '8px 1fr 32px',
        sm:   '8px 1fr 96px 32px',
        md:   '8px 1fr 160px 96px 32px',
        lg:   '8px 1fr 160px 160px 88px 88px 44px',
      }}
      alignItems="center"
      columnGap={4}
      px={4}
      py={3}
      bg="cardBg"
      borderRadius="xl"
      border="1px solid"
      borderColor="cardBorder"
      boxShadow="0 1px 3px rgba(15,23,42,0.03)"
      _hover={{
        borderColor: 'rgba(15,23,42,0.10)',
        boxShadow: '0 2px 8px rgba(15,23,42,0.06)',
        bg: 'subtleBg',
      }}
      transition="all 0.15s"
      minH="60px"
    >
      {/* Col 1: Dot — gray when completed */}
      <Box
        w="8px" h="8px" borderRadius="full" flexShrink={0}
        {...(isCompleted
          ? { bg: 'gray.400', _dark: { bg: 'gray.500' } }
          : { bgGradient: gradient }
        )}
      />

      {/* Col 2: Title + category badges — dimmed when completed */}
      <Box minW={0} opacity={isCompleted ? 0.6 : 1}>
        <HStack mb={0.5} spacing={1.5} flexWrap="nowrap" overflow="hidden">
          {isCompleted && (
            <Badge colorScheme="gray" variant="solid" fontSize="9px" borderRadius="md" flexShrink={0}>
              {t('eventCard.statusCompleted')}
            </Badge>
          )}
          <Badge
            colorScheme={event.category === 'business' ? 'gray' : 'green'}
            variant="subtle"
            fontSize="9px"
            borderRadius="md"
            flexShrink={0}
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
            flexShrink={0}
          >
            <Icon as={event.is_public ? FiGlobe : FiLock} boxSize={2.5} />
            {event.is_public ? t('eventCard.publicBadge') : t('eventCard.personalBadge')}
          </Badge>
        </HStack>
        <Text fontSize="sm" fontWeight="600" color="mainText" noOfLines={1}>
          {event.title}
        </Text>
      </Box>

      {/* Col 3: Date (sm+) */}
      <HStack spacing={1.5} minW={0} display={{ base: 'none', sm: 'flex' }}>
        <Icon as={FiCalendar} color="brand.400" boxSize={3.5} flexShrink={0} />
        <Text fontSize="xs" color="dimText" fontWeight="500" noOfLines={1}>
          {formatDate(event.date_start)}, {formatTime(event.date_start)}
        </Text>
      </HStack>

      {/* Col 4: Location (md+) */}
      <HStack spacing={1.5} minW={0} display={{ base: 'none', md: 'flex' }}>
        <Icon as={FiMapPin} color="brand.400" boxSize={3.5} flexShrink={0} />
        <Text fontSize="xs" color="dimText" fontWeight="500" noOfLines={1}>
          {event.location}
        </Text>
      </HStack>

      {/* Col 5: Status badge (lg+) — always renders to hold the column */}
      <Box display={{ base: 'none', lg: 'flex' }} justifyContent="flex-start">
        {colors && statusKey ? (
          <Box px={2.5} py={1} borderRadius="full" bg={colors.bg}>
            <Text fontSize="xs" fontWeight="700" color={colors.color} whiteSpace="nowrap">
              {t(`eventCard.status${statusKey.charAt(0).toUpperCase()}${statusKey.slice(1)}`)}
            </Text>
          </Box>
        ) : null}
      </Box>

      {/* Col 6: Participants (lg+) */}
      <HStack spacing={1.5} display={{ base: 'none', lg: 'flex' }}>
        <Text fontSize="xs" color="faintText">{confirmedCount}</Text>
        <AvatarGroup size="xs" max={3} spacing="-6px">
          {event.participants
            .filter(p => p.status === 'confirmed')
            .map(p => (
              <Avatar key={p.id} name={p.user.name} bg="brand.400" color="white" />
            ))}
        </AvatarGroup>
      </HStack>

      {/* Col 7: Action */}
      <Box display="flex" justifyContent="flex-end">
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
          <Icon as={FiChevronRight} color="faintText" boxSize={4} />
        )}
      </Box>
    </Box>
  )
}
