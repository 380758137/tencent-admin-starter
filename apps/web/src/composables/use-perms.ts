import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'

export function usePerm(perm: string) {
  const auth = useAuthStore()
  return computed(() => auth.hasPerm(perm))
}

export function useCrudPerms(prefix: string) {
  const auth = useAuthStore()
  return {
    canCreate: computed(() => auth.hasPerm(`${prefix}:create`)),
    canUpdate: computed(() => auth.hasPerm(`${prefix}:update`)),
    canDelete: computed(() => auth.hasPerm(`${prefix}:delete`))
  }
}
