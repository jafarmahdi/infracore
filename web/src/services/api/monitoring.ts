import type {
  CreateHostRequest,
  HostStatusCounts,
  ListHostsResponse,
  MonitoredHost,
} from "@/types/monitoring";
import { apiClient } from "./client";

export interface ListHostsParams {
  page?: number;
  page_size?: number;
  status?: string;
  monitoring_type?: string;
  search?: string;
}

export async function listHosts(params?: ListHostsParams): Promise<ListHostsResponse> {
  const { data } = await apiClient.get<ListHostsResponse>("/monitoring/hosts", { params });
  return data;
}

export async function getHost(id: string): Promise<MonitoredHost> {
  const { data } = await apiClient.get<MonitoredHost>(`/monitoring/hosts/${id}`);
  return data;
}

export async function createHost(body: CreateHostRequest): Promise<MonitoredHost> {
  const { data } = await apiClient.post<MonitoredHost>("/monitoring/hosts", body);
  return data;
}

export async function deleteHost(id: string): Promise<void> {
  await apiClient.delete(`/monitoring/hosts/${id}`);
}

export async function getStatusCounts(): Promise<HostStatusCounts> {
  const { data } = await apiClient.get<HostStatusCounts>("/monitoring/hosts/counts");
  return data;
}
