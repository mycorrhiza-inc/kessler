"use client"
import React, { useState, useCallback, useMemo, useEffect } from 'react';
import { FilterMultiSelectRaw, FilterFieldDefinition, Option } from './FilterMultiSelectRaw';

export interface DynamicMultiSelectProps {
  fieldDefinition: FilterFieldDefinition;
  value: string;
  onChange: (value: string) => void;
  onFocus?: () => void;
  onBlur?: () => void;
  disabled?: boolean;
  className?: string;
}

export function DynamicMultiSelect({
  fieldDefinition,
  value,
  onChange,
  onFocus,
  onBlur,
  disabled = false,
  className = "",
}: DynamicMultiSelectProps) {
  const [selectedValues, setSelectedValues] = useState<string[]>(
    value ? value.split(',').filter(Boolean) : []
  );
  const [isOpen, setIsOpen] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  const handleSelectionChange = useCallback((selectedOptions: string[]) => {
    setSelectedValues(selectedOptions);
    onChange(selectedOptions.join(','));
  }, [onChange]);

  const toggleOption = useCallback((optionValue: string) => {
    const newSelection = selectedValues.includes(optionValue)
      ? selectedValues.filter(v => v !== optionValue)
      : [...selectedValues, optionValue];
    handleSelectionChange(newSelection);
  }, [selectedValues, handleSelectionChange]);

  const removeSelectedItem = useCallback((optionValue: string) => {
    const newSelection = selectedValues.filter(v => v !== optionValue);
    handleSelectionChange(newSelection);
  }, [selectedValues, handleSelectionChange]);

  const filteredOptions: Option[] = useMemo(() => {
    if (!searchTerm) return fieldDefinition.options || [];
    return (fieldDefinition.options || []).filter(option =>
      option.label.toLowerCase().includes(searchTerm.toLowerCase()) ||
      option.value.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [fieldDefinition.options, searchTerm]);

  const selectedLabels = useMemo(() => {
    return selectedValues.map(val => {
      const opt = fieldDefinition.options?.find(o => o.value === val);
      return opt ? opt.label : val;
    });
  }, [selectedValues, fieldDefinition.options]);

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

  const clearAll = useCallback(() => {
    handleSelectionChange([]);
  }, [handleSelectionChange]);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const target = event.target as Element;
      if (isOpen && !target.closest('.multiselect-container')) {
        setIsOpen(false);
        onBlur?.();
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen, onBlur]);

  return (
    <FilterMultiSelectRaw
      fieldDefinition={fieldDefinition}
      selectedValues={selectedValues}
      selectedLabels={selectedLabels}
      filteredOptions={filteredOptions}
      isOpen={isOpen}
      searchTerm={searchTerm}
      disabled={disabled}
      className={className}
      onToggleDropdown={handleDropdownToggle}
      onSearchChange={handleSearchChange}
      onClearSearch={clearSearch}
      onToggleOption={toggleOption}
      onRemoveSelectedItem={removeSelectedItem}
      onClearAll={clearAll}
    />
  );
}