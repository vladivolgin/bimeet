import {
  AlertDialog,
  AlertDialogBody,
  AlertDialogContent,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogOverlay,
  Badge,
  Box,
  Button,
  Collapse,
  Divider,
  Heading,
  HStack,
  Icon,
  IconButton,
  Menu,
  MenuButton,
  MenuDivider,
  MenuItem,
  MenuList,
  Text,
  Tooltip,
  useDisclosure,
  useToast,
  VStack,
} from '@chakra-ui/react'
import { useNavigate } from 'react-router-dom'
import { useRef } from 'react'
import { useTranslation } from 'react-i18next'
import {
  FiBriefcase,
  FiCalendar,
  FiChevronDown,
  FiChevronUp,
  FiClock,
  FiGlobe,
  FiLock,
  FiMapPin,
  FiMoreVertical,
  FiTag,
  FiEdit2,
  FiTrash2,
  FiUsers,
} from 'react-icons/fi'
import type { Event } from '@/entities/event/model/types'
import { getGradient } from '@/widgets/event-card'
import { formatDate, formatTime } from '@/shared/lib/formatDate'
import { useDeleteEvent, useJoinPublicEvent } from '@/features/event-manage/model/hooks'
import { useAuth } from '@/features/auth/model/AuthContext'

interface EventInfoPanelProps {
  event: Event
  isOrganizer: boolean
}

export default function EventInfoPanel({ event, isOrganizer }: EventInfoPanelProps) {
  const navigate = useNavigate()
  const { t } = useTranslation()
  const toast = useToast()
  const { user } = useAuth()
  const { isOpen: isDescExpanded, onToggle: onDescToggle } = useDisclosure()
  const {
    isOpen: isDeleteOpen, onOpen: onDeleteOpen, onClose: onDeleteClose,
  } = useDisclosure()
  const cancelDeleteRef = useRef<HTMLButtonElement>(null)

  const deleteEvent = useDeleteEvent()
  const joinPublicEvent = useJoinPublicEvent()

  const isParticipant = event.participants.some(p => p.user?.id === user?.id)
  const confirmedCount = event.confirmed_count ?? event.participants.filter(p => p.status === 'confirmed').length
  const spotsLeft = event.max_guests != null ? event.max_guests - confirmedCount : null
  const canJoin = event.is_public && event.status === 'active' && !isOrganizer && !isParticipant && (spotsLeft == null || spotsLeft > 0)

  function handleJoin() {
    joinPublicEvent.mutate(event.id, {
      onSuccess: () => toast({ title: t('eventInfo.joinSuccess'), status: 'success', duration: 2000 }),
      onError: (err) => toast({ title: t('common.error'), description: err.message, status: 'error', duration: 3000 }),
    })
  }

  const gradient = getGradient(event.id, event.category)
  const sameDay = formatDate(event.date_start) === formatDate(event.date_end)
  const descLong = (event.description ?? '').length > 120

  function handleDelete() {
    deleteEvent.mutate(event.id, {
      onSuccess: () => {
        onDeleteClose()
        navigate('/events')
      },
    })
  }

  return (
    <>
    <Box
      bg="cardBg"
      borderRadius="xl"
      border="1px solid"
      borderColor="cardBorder"
      boxShadow="0 1px 3px rgba(15,23,42,0.04), 0 4px 16px rgba(15,23,42,0.05)"
    >
      <Box h="4px" bgGradient={gradient} borderTopRadius="xl" />

      <Box p={5}>
        <HStack justify="space-between" align="flex-start" mb={3}>
          <Box flex={1} minW={0}>
            <HStack mb={1} spacing={2} flexWrap="wrap">
              <Text
                fontSize="10px"
                fontWeight="600"
                color={event.category === 'business' ? 'brand.700' : 'brand.500'}
                textTransform="uppercase"
                letterSpacing="0.1em"
              >
                {event.category === 'business' ? t('eventInfo.categoryBusiness') : t('eventInfo.categoryOrdinary')}
              </Text>
              {event.status === 'completed' && (
                <Badge colorScheme="gray" variant="subtle" fontSize="9px" borderRadius="md">
                  {t('eventInfo.statusCompleted')}
                </Badge>
              )}
            </HStack>
            <Heading size="md" letterSpacing="-0.3px" lineHeight="1.3">
              {event.title}
            </Heading>
          </Box>

          {isOrganizer && (
            <Menu>
              <MenuButton
                as={IconButton}
                icon={<FiMoreVertical />}
                variant="ghost"
                size="sm"
                aria-label={t('eventInfo.menuActions')}
                borderRadius="xl"
                color="faintText"
                _hover={{ bg: 'subtleBg', color: 'dimText' }}
                flexShrink={0}
              />
              <MenuList>
                {event.status === 'active' && (
                  <>
                    <MenuItem fontSize="sm" icon={<FiEdit2 size={14} />} onClick={() => navigate(`/events/${event.id}/edit`)}>
                      {t('common.edit')}
                    </MenuItem>
                    <MenuDivider />
                  </>
                )}
                <MenuItem fontSize="sm" icon={<FiTrash2 size={14} />} color="red.400" onClick={onDeleteOpen}>
                  {t('common.delete')}
                </MenuItem>
              </MenuList>
            </Menu>
          )}
        </HStack>

        {event.description && (
          <Box mb={4}>
            <Collapse in={isDescExpanded} startingHeight={descLong ? '3.6em' : undefined}>
              <Text color="dimText" fontSize="sm" lineHeight="1.6">
                {event.description}
              </Text>
            </Collapse>
            {descLong && (
              <Text
                as="button"
                fontSize="xs"
                color="brand.500"
                fontWeight="500"
                mt={1}
                display="flex"
                alignItems="center"
                gap={0.5}
                onClick={onDescToggle}
                _hover={{ color: 'brand.600' }}
              >
                {isDescExpanded ? (
                  <><Icon as={FiChevronUp} /> {t('eventInfo.descCollapse')}</>
                ) : (
                  <><Icon as={FiChevronDown} /> {t('eventInfo.descExpand')}</>
                )}
              </Text>
            )}
          </Box>
        )}

        <VStack align="stretch" spacing={2}>
          {sameDay ? (
            <>
              <HStack spacing={2.5}>
                <Icon as={FiCalendar} color="brand.400" boxSize={3.5} flexShrink={0} />
                <Text fontSize="sm" color="dimText" fontWeight="500">{formatDate(event.date_start)}</Text>
              </HStack>
              <HStack spacing={2.5}>
                <Icon as={FiClock} color="brand.400" boxSize={3.5} flexShrink={0} />
                <Text fontSize="sm" color="dimText" fontWeight="500">
                  {formatTime(event.date_start)} — {formatTime(event.date_end)}
                </Text>
              </HStack>
            </>
          ) : (
            <HStack spacing={2.5}>
              <Icon as={FiCalendar} color="brand.400" boxSize={3.5} flexShrink={0} />
              <Text fontSize="sm" color="dimText" fontWeight="500">
                {formatDate(event.date_start)}, {formatTime(event.date_start)} — {formatDate(event.date_end)}, {formatTime(event.date_end)}
              </Text>
            </HStack>
          )}

          <HStack spacing={2.5}>
            <Icon as={FiMapPin} color="brand.400" boxSize={3.5} flexShrink={0} />
            <Text fontSize="sm" color="dimText" fontWeight="500">{event.location}</Text>
          </HStack>

          {event.dress_code && (
            <HStack spacing={2.5}>
              <Icon as={FiTag} color="brand.400" boxSize={3.5} flexShrink={0} />
              <Text fontSize="sm" color="dimText" fontWeight="500">{t('eventInfo.dressCode', { value: event.dress_code })}</Text>
            </HStack>
          )}

          {event.category === 'business' && (
            <HStack spacing={2.5}>
              <Icon as={FiBriefcase} color="brand.400" boxSize={3.5} flexShrink={0} />
              <Badge colorScheme="blue" variant="subtle" borderRadius="md" fontSize="xs">
                {t('eventInfo.categoryBusiness')}
              </Badge>
            </HStack>
          )}

          <HStack spacing={2.5}>
            <Icon as={event.is_public ? FiGlobe : FiLock} color="brand.400" boxSize={3.5} flexShrink={0} />
            <Badge colorScheme={event.is_public ? 'teal' : 'purple'} variant="subtle" borderRadius="md" fontSize="xs">
              {event.is_public ? t('eventCard.publicBadge') : t('eventCard.personalBadge')}
            </Badge>
          </HStack>

          <HStack spacing={2.5}>
            <Icon as={FiUsers} color="brand.400" boxSize={3.5} flexShrink={0} />
            <Text fontSize="sm" color="dimText" fontWeight="500">
              {event.max_guests != null
                ? `${confirmedCount} / ${event.max_guests}`
                : t('eventInfo.confirmedCount', { count: confirmedCount })}
            </Text>
          </HStack>
        </VStack>

        {(canJoin || (event.is_public && !isOrganizer && !isParticipant && spotsLeft != null && spotsLeft <= 0)) && (
          <Box mt={4}>
            <Tooltip label={!canJoin ? t('eventInfo.noSpots') : undefined} isDisabled={canJoin}>
              <Button
                w="full"
                colorScheme="blue"
                isDisabled={!canJoin}
                isLoading={joinPublicEvent.isPending}
                onClick={handleJoin}
              >
                {t('eventInfo.join')}
              </Button>
            </Tooltip>
          </Box>
        )}

        <Divider my={4} borderColor="subtleBorder" />

        <HStack spacing={2}>
          <Text fontSize="xs" color="faintText">{t('eventInfo.organizer')}</Text>
          <Text fontSize="xs" fontWeight="600" color="dimText">{event.organizer.name}</Text>
        </HStack>
      </Box>
    </Box>

    {/* Delete confirmation */}
    <AlertDialog isOpen={isDeleteOpen} leastDestructiveRef={cancelDeleteRef} onClose={onDeleteClose}>
      <AlertDialogOverlay>
        <AlertDialogContent>
          <AlertDialogHeader fontSize="lg" fontWeight="bold">
            {t('eventInfo.deleteDialogTitle')}
          </AlertDialogHeader>
          <AlertDialogBody>
            {t('eventInfo.deleteDialogBody')}
          </AlertDialogBody>
          <AlertDialogFooter>
            <Button ref={cancelDeleteRef} onClick={onDeleteClose} variant="ghost">
              {t('common.cancel')}
            </Button>
            <Button
              colorScheme="red"
              ml={3}
              onClick={handleDelete}
              isLoading={deleteEvent.isPending}
            >
              {t('eventInfo.deleteDialogConfirm')}
            </Button>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogOverlay>
    </AlertDialog>
    </>
  )
}
