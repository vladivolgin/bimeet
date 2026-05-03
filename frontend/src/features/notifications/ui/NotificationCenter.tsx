import {
  Badge,
  Box,
  HStack,
  IconButton,
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  Text,
  Tooltip,
  VStack,
  useColorModeValue,
} from '@chakra-ui/react'
import { FiBell, FiCheck, FiTrash2 } from 'react-icons/fi'
import { useTranslation } from 'react-i18next'
import { useNavigate } from 'react-router-dom'
import { useNotifications } from '@/entities/notification/queries'
import {
  useDeleteAllNotifications,
  useDeleteNotification,
  useMarkAllRead,
  useMarkRead,
} from '@/features/notifications/model/hooks'
import { formatDate } from '@/shared/lib/formatDate'
import type { Notification } from '@/entities/notification/model/types'

function notificationHref(n: Notification): string | null {
  if (!n.event_id) return null
  if (n.type === 'collection_contribution_pending') return `/events/${n.event_id}?tab=collections`
  return `/events/${n.event_id}`
}

export default function NotificationCenter() {
  const { data: notifications = [] } = useNotifications()
  const markRead = useMarkRead()
  const markAllRead = useMarkAllRead()
  const deleteOne = useDeleteNotification()
  const deleteAll = useDeleteAllNotifications()
  const { t } = useTranslation()
  const navigate = useNavigate()

  const unreadCount = notifications.filter((n) => !n.is_read).length
  const triggerHoverBg = useColorModeValue('rgba(15,23,42,0.05)', 'rgba(255,255,255,0.07)')
  const dotColor = useColorModeValue('blue.500', 'blue.300')

  return (
    <Popover placement="bottom-end" isLazy>
      <PopoverTrigger>
        <Box position="relative" display="inline-flex">
          <Tooltip label={t('notifications.title')} borderRadius="lg">
            <IconButton
              aria-label={t('notifications.title')}
              icon={<FiBell size={16} />}
              variant="ghost"
              size="sm"
              borderRadius="10px"
              color="dimText"
              _hover={{ color: 'mainText', bg: triggerHoverBg }}
            />
          </Tooltip>
          {unreadCount > 0 && (
            <Badge
              position="absolute"
              top="-2px"
              right="-2px"
              colorScheme="blue"
              borderRadius="full"
              fontSize="9px"
              minW="16px"
              h="16px"
              display="flex"
              alignItems="center"
              justifyContent="center"
              px={1}
            >
              {unreadCount > 99 ? '99+' : unreadCount}
            </Badge>
          )}
        </Box>
      </PopoverTrigger>

      <PopoverContent
        w="360px"
        maxH="480px"
        overflow="hidden"
        display="flex"
        flexDirection="column"
        bg="cardBg"
        borderColor="defaultBorder"
        boxShadow="0 8px 32px rgba(0,0,0,0.12)"
        _dark={{ boxShadow: '0 8px 32px rgba(0,0,0,0.45)' }}
      >
        <PopoverHeader px={4} py={3} borderBottomWidth="1px" borderColor="subtleBorder">
          <HStack justify="space-between">
            <Text fontWeight="600" fontSize="sm" color="mainText">
              {t('notifications.title')}
            </Text>
            <HStack spacing={1}>
              <Tooltip label={t('notifications.markAllRead')} borderRadius="lg" openDelay={400}>
                <IconButton
                  aria-label={t('notifications.markAllRead')}
                  icon={<FiCheck size={14} />}
                  size="xs"
                  variant="ghost"
                  color="dimText"
                  borderRadius="8px"
                  _hover={{ color: 'green.500', bg: 'subtleBg' }}
                  onClick={() => markAllRead.mutate()}
                  isLoading={markAllRead.isPending}
                  isDisabled={unreadCount === 0}
                />
              </Tooltip>
              <Tooltip label={t('notifications.deleteAll')} borderRadius="lg" openDelay={400}>
                <IconButton
                  aria-label={t('notifications.deleteAll')}
                  icon={<FiTrash2 size={14} />}
                  size="xs"
                  variant="ghost"
                  color="dimText"
                  borderRadius="8px"
                  _hover={{ color: 'red.400', bg: 'subtleBg' }}
                  onClick={() => deleteAll.mutate()}
                  isLoading={deleteAll.isPending}
                  isDisabled={notifications.length === 0}
                />
              </Tooltip>
            </HStack>
          </HStack>
        </PopoverHeader>

        <PopoverBody p={0} overflowY="auto">
          {notifications.length === 0 ? (
            <Box py={8} textAlign="center">
              <Text fontSize="sm" color="dimText">
                {t('notifications.empty')}
              </Text>
            </Box>
          ) : (
            <VStack spacing={0} align="stretch">
              {notifications.map((n) => {
                const href = notificationHref(n)
                return (
                <HStack
                  key={n.id}
                  px={4}
                  py={3}
                  spacing={3}
                  bg="cardBg"
                  borderBottomWidth="1px"
                  borderColor="subtleBorder"
                  align="flex-start"
                  role="group"
                  cursor={href ? 'pointer' : 'default'}
                  _hover={{ bg: 'subtleBg' }}
                  transition="background 0.15s"
                  onClick={href ? () => {
                    if (!n.is_read) markRead.mutate(n.id)
                    navigate(href)
                  } : undefined}
                >
                  <Box pt="5px" flexShrink={0}>
                    <Box
                      w="7px"
                      h="7px"
                      borderRadius="full"
                      bg={n.is_read ? 'transparent' : dotColor}
                      border="1.5px solid"
                      borderColor={n.is_read ? 'defaultBorder' : dotColor}
                    />
                  </Box>

                  <Box flex={1} minW={0}>
                    <Text
                      fontSize="sm"
                      color={n.is_read ? 'dimText' : 'mainText'}
                      fontWeight={n.is_read ? 400 : 500}
                      noOfLines={2}
                    >
                      {n.message}
                    </Text>
                    <Text fontSize="xs" color="faintText" mt={0.5}>
                      {formatDate(n.created_at)}
                    </Text>
                  </Box>

                  <HStack
                    spacing={0.5}
                    flexShrink={0}
                    opacity={0}
                    _groupHover={{ opacity: 1 }}
                    transition="opacity 0.15s"
                  >
                    {!n.is_read && (
                      <IconButton
                        aria-label={t('notifications.markRead')}
                        icon={<FiCheck size={13} />}
                        size="xs"
                        variant="ghost"
                        color="dimText"
                        borderRadius="8px"
                        _hover={{ color: 'green.500', bg: 'transparent' }}
                        onClick={(e) => { e.stopPropagation(); markRead.mutate(n.id) }}
                      />
                    )}
                    <IconButton
                      aria-label={t('notifications.deleteOne')}
                      icon={<FiTrash2 size={13} />}
                      size="xs"
                      variant="ghost"
                      color="dimText"
                      borderRadius="8px"
                      _hover={{ color: 'red.400', bg: 'transparent' }}
                      onClick={(e) => { e.stopPropagation(); deleteOne.mutate(n.id) }}
                    />
                  </HStack>
                </HStack>
              )})}
            </VStack>
          )}
        </PopoverBody>
      </PopoverContent>
    </Popover>
  )
}
