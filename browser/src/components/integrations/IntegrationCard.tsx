import type { IntegrationPayload } from "@/types";
import { GitHubCard } from "./GitHubCard";
import { GenericIntegrationCard } from "./GenericIntegrationCard";

interface IntegrationCardProps {
  name: string;
  payload: IntegrationPayload;
  onSync: () => void;
  isSyncing: boolean;
}

const INTEGRATION_COMPONENTS: Record<
  string,
  React.ComponentType<{ payload: IntegrationPayload; onSync: () => void; isSyncing: boolean }>
> = {
  github: GitHubCard,
};

export function IntegrationCard({ name, payload, onSync, isSyncing }: IntegrationCardProps) {
  const Custom = INTEGRATION_COMPONENTS[name];
  if (Custom) return <Custom payload={payload} onSync={onSync} isSyncing={isSyncing} />;
  return <GenericIntegrationCard name={name} payload={payload} onSync={onSync} isSyncing={isSyncing} />;
}
