"use client"
import React, { useState, useCallback, useMemo, useEffect } from 'react';
import {
  FilterSingleSelectRaw,
  FilterFieldDefinition,
  Option,
} from './FilterSingleSelectRaw';

export interface DynamicSingleSelectProps {
  fieldDefinition: FilterFieldDefinition;
  value: string;
  onChange: (value: string) => void;
  onFocus?: () => void;
  onBlur?: () => void;
  disabled?: boolean;
  className?: string;
  allowClear?: boolean;
  defaultValue?: string;
  dynamicWidth?: boolean;
  minWidth?: string;
  maxWidth?: string;
}

export function DynamicSingleSelect({
  fieldDefinition,
  value,
  onChange,
  onFocus,
  onBlur,
  disabled = false,
  className = '',
  allowClear = true,
  defaultValue = '',
  dynamicWidth = false,
  minWidth = '200px',
  maxWidth = '150px',
}: DynamicSingleSelectProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [componentWidth, setComponentWidth] = useState<string>('100%');

  useEffect(() => {
    if (!value && defaultValue &&
        fieldDefinition.options?.some(opt => opt.value === defaultValue)
    ) {
      onChange(defaultValue);
    }
  }, [value, defaultValue, onChange, fieldDefinition.options]);

  const handleSelectionChange = useCallback((optValue: string) => {
    onChange(optValue);
    setIsOpen(false);
  }, [onChange]);

  const clearSelection = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onChange('');
  }, [onChange]);

  const filteredOptions: Option[] = useMemo(() => {
    if (!searchTerm) return fieldDefinition.options || [];
    return (fieldDefinition.options || []).filter(option =>
      option.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
      option.value.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [fieldDefinition.options, searchTerm]);

  const selectedOption: Option | undefined = useMemo(() => {
    return fieldDefinition.options?.find(opt => opt.value === value);
  }, [value, fieldDefinition.options]);

  const calculateDynamicWidth = useCallback(() => {
    if (!dynamicWidth) return '100%';
    let calculated = '200px';
    if (selectedOption) {
      const totalLen = selectedOption.label.length + fieldDefinition.displayName.length + 10;
      const estWidth = Math.max(200, totalLen * 8 + 80);
      calculated = `${Math.min(parseInt(maxWidth), estWidth)}px`;
    } else if (fieldDefinition.placeholder) {
      const plLen = fieldDefinition.placeholder.length;
      const estWidth = Math.max(200, plLen * 8 + 40);
      calculated = `${Math.min(parseInt(maxWidth), estWidth)}px`;
    }
    return calculated;
  }, [selectedOption, fieldDefinition, dynamicWidth, maxWidth]);

  useEffect(() => {
    if (dynamicWidth) {
      setComponentWidth(calculateDynamicWidth());
    }
  }, [selectedOption, dynamicWidth, calculateDynamicWidth]);

  const handleDropdownToggle = useCallback(() => {
    if (!disabled) {
      setIsOpen(prev => !prev);
      if (!isOpen) onFocus?.(); else onBlur?.();
    }
  }, [disabled, isOpen, onFocus, onBlur]);

  const handleSearchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  }, []);

  const clearSearch = useCallback(() => {
    setSearchTerm('');
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (isOpen && !target.closest('.singleselect-container')) {
        setIsOpen(false);
        onBlur?.();
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen, onBlur]);

  const containerStyle = dynamicWidth
    ? { width: componentWidth, minWidth, maxWidth, transition: 'width 0.2s ease-in-out' }
    : undefined;

  return (
    <FilterSingleSelectRaw
      fieldDefinition={fieldDefinition}
      value={value}
      selectedOption={selectedOption}
      filteredOptions={filteredOptions}
      isOpen={isOpen}
      searchTerm={searchTerm}
      disabled={disabled}
      className={className}
      allowClear={allowClear}
      containerStyle={containerStyle}
      onToggleDropdown={handleDropdownToggle}
      onSearchChange={handleSearchChange}
      onClearSearch={clearSearch}
      onSelectOption={handleSelectionChange}
      onClearSelection={() => handleSelectionChange('')}
    />
  );
}